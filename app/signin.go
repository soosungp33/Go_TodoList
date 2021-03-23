package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// https://pkg.go.dev/golang.org/x/oauth2#pkg-overview

// 로그인 정보를 나타내는 구조체
type GoogleUserId struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:verified_email"`
	Picture       string `json:"picture"`
}

var googleOauthConfig = oauth2.Config{
	// RedirectURL:  "http://localhost:3000/auth/google/callback", -> 로컬 환경
	RedirectURL:  os.Getenv("DOMAIN_NAME") + "/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"), // 환경변수에 설정해놓은 ID를 가져옴
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)                   // 상태쿠키를 만든다.(구글에서 알려준 state와 저장한 쿠키가 같으면 공격에 의한 것이 아님을 판단)
	url := googleOauthConfig.AuthCodeURL(state)            // 유저를 어떤 경로로 보내야하는지 알려준다.
	http.Redirect(w, r, url, http.StatusTemporaryRedirect) // 해당 경로로 리다이렉트된다.(콜백 url)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour) // 하루 동안 유지

	b := make([]byte, 16)
	rand.Read(b) // 바이트를 랜덤하게 채운다.
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration} // 랜덤한 상태값으로 쿠키를 만들어준다.(이름과 값은 필수)
	http.SetCookie(w, cookie)

	return state
}

func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate") // 위에서 저장했던 쿠키를 반환
	// r.FormValue("state") -> 구글에서 보내준 state
	if r.FormValue("state") != oauthstate.Value { // 둘이 다르면 잘못된 접근
		errMsg := fmt.Sprintf("invalid google oauth state:%s state:%s", oauthstate.Value, r.FormValue("state"))
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
	}
	data, err := getGoogleUserInfo(r.FormValue("code")) // userinfo를 구글에 request해서 받아온다.

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// data에는 id, email, verified_email, picture 가 나타난다.
	// ID 정보를 세션 쿠키에다 저장(다른 페이지에서 로그인 여부를 판단하는데 사용)
	var userInfo GoogleUserId
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 참고 : https://github.com/gorilla/sessions
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set some session values.
	session.Values["id"] = userInfo.ID
	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w) // session이라는 이름의 세션에 저장
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // 최종적으로 로그인을 한 뒤 정보들을 세션에 저장하고 메인으로 리다이렉션
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body)
}
