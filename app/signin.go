package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
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
	googlestate := r.FormValue("state")     // 구글에서 보내준 state
	if googlestate != oauthstate.Value {    // 둘이 다르면 잘못된 접근
		log.Printf("invalid google oauth state:%s state:%s", oauthstate.Value, googlestate)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // 잘못된 접근이면 루트경로로 리다이렉트
	}
	data, err := getGoogleUserInfo(r.FormValue("code")) // userinfo를 구글에 request해서 받아온다.
	if err != nil {
		log.Println(err.Error)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprint(w, string(data)) // string으로 화면에 내 정보가 나타난다.
	/*
		id, email, verified_email, picture 가 나타난다.
		id를 통해서 키를 만들면 됨
	*/
}

var oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

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
