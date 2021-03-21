package app

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/soosungp33/Go_TodoList/model"
	"github.com/unrolled/render"
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
	list := a.db.GetTodos()
	// 렌더링을 사용해서 JSON으로 반환
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) { // 프론트에서 올 때 name에다 Item을 넣어서 온다.
	name := r.FormValue("name") // 키를 통해서 value에 있는 Item을 꺼낸다.
	todo := a.db.AddTodo(name)
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

func MakeHandler(filepath string) *AppHandler {
	r := mux.NewRouter()
	a := &AppHandler{
		Handler: r,
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
