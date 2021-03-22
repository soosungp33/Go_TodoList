package app

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/soosungp33/Go_TodoList/model"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func getSessionId(r *http.Request) string { // indexHandler가 호출될 때 세션 ID를 읽어와야 로그인 정보들을 활용할 수 있다.(쿠키는 Request에 들어있다.)
	session, err := store.Get(r, "session") // session이라는 세션에 저장되어 있는 정보에 접근
	if err != nil {
		return ""
	}

	val := session.Values["id"]
	if val == nil { // id가 없으면(로그인을 안했으면)
		return ""
	}

	return val.(string)
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY"))) // 환경 변수에 설정해놓은 SESSION_KEY를 가져옴(임의로 설정하면 됨)
var rd *render.Render = render.New()                                  // 렌더링을 사용해서 JSON 변환을 쉽게하기

type AppHandler struct {
	http.Handler
	db model.DBHandler
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionId(r) // 세션 아이디에 해당하는 리스트를 가져옴
	list := a.db.GetTodos(sessionId)
	// 렌더링을 사용해서 JSON으로 반환
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) { // 프론트에서 올 때 name에다 Item을 넣어서 온다.
	sessionId := getSessionId(r) // 세션 아이디에 해당하는 리스트에 add
	name := r.FormValue("name")  // 키를 통해서 value에 있는 Item을 꺼낸다.
	todo := a.db.AddTodo(sessionId, name)
	// 렌더링을 사용해서 JSON으로 반환
	rd.JSON(w, http.StatusCreated, todo)
}

func (a *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)               // mux.Vars를 통해 변수 맵을 만들어 검색
	id, _ := strconv.Atoi(vars["id"]) // id에 있는 문자열을 숫자로 변경
	ok := a.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func (a *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                           // mux.Vars를 통해 변수 맵을 만들어 검색
	id, _ := strconv.Atoi(vars["id"])             // id에 있는 문자열을 숫자로 변경
	complete := r.FormValue("complete") == "true" // complete에 true가 담겨오면 체크하라는 뜻이고, false면 체크를 해제하라는 뜻
	ok := a.db.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

type Success struct {
	Success bool `json:"success"`
}

func (a *AppHandler) Close() {
	a.db.Close()
}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 유저가 요청한 URL이 로그인 관련이면 그냥 next해야된다(안해주면 무한루프를 돌 수도 있음)
	if strings.Contains(r.URL.Path, "/signin") || strings.Contains(r.URL.Path, "/auth") {
		// URL에 로그인 페이지가 포함되어 있거나 로그인 버튼을 눌러서 회원가입(구글 권한 요청 페이지) 페이지가 포함되어 있으면 next 핸들러로 간다.
		next(w, r)
		return
	}

	// 로그인 URL외에 다른 URL을 요청할 때 유저가 로그인 되어있으면 next로 아니면 로그인으로 리다이렉트
	sessionID := getSessionId(r)
	if sessionID != "" {
		next(w, r)
		return
	}

	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)

}

func MakeHandler(filepath string) *AppHandler {
	r := mux.NewRouter()
	// 미들웨어를 추가해서 모든 핸들러가 불릴 때마다 세션아이디를 체크
	// 원래 negroni.Classic()은 NewRecovery(), NewLogger(), NewStatic(http.Dir("public"))을 순서대로 반환한다.
	// public으로 가기전에 CheckSignin핸들러를 불러 로그인이 되어있는지 아닌지 검사한다.
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(CheckSignin),
		negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	a := &AppHandler{
		Handler: n,
		db:      model.NewDBHandler(filepath),
	}

	r.HandleFunc("/todos", a.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", a.addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", a.removeTodoHandler).Methods("DELETE") // 아이디는 1개 이상인 숫자로만 이루어짐
	r.HandleFunc("/complete-todo/{id:[0-9]+}", a.completeTodoHandler).Methods("GET")

	// 로그인 부분
	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/auth/google/callback", googleAuthCallback)

	// 처음 서버를 시작하고 local:3000으로 들어갔을 때 바로 local:3000/todo.html로 리다이렉트 하는 핸들러
	r.HandleFunc("/", indexHandler)

	return a
}
