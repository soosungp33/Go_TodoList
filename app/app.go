package app

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/soosungp33/Go_TodoList/model"
	"github.com/unrolled/render"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

var rd *render.Render = render.New() // 렌더링을 사용해서 JSON 변환을 쉽게하기

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

func MakeHandler() *AppHandler {
	r := mux.NewRouter()
	a := &AppHandler{
		Handler: r,
		db:      model.NewDBHandler(),
	}

	r.HandleFunc("/todos", a.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", a.addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", a.removeTodoHandler).Methods("DELETE") // 아이디는 1개 이상인 숫자로만 이루어짐
	r.HandleFunc("/complete-todo/{id:[0-9]+}", a.completeTodoHandler).Methods("GET")
	// 처음 서버를 시작하고 local:3000으로 들어갔을 때 바로 local:3000/todo.html로 리다이렉트 하는 핸들러
	r.HandleFunc("/", indexHandler)

	return a
}
