package app

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

var rd *render.Render

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var todoMap map[int]*Todo // 인메모리 db역할

func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	list := []*Todo{}
	for _, v := range todoMap {
		list = append(list, v)
	}
	// 렌더링을 사용해서 JSON으로 반환
	rd.JSON(w, http.StatusOK, list)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) { // 프론트에서 올 때 name에다 Item을 넣어서 온다.
	name := r.FormValue("name") // 키를 통해서 value에 있는 Item을 꺼낸다.
	id := len(todoMap) + 1      // 임의의 ID
	todo := &Todo{id, name, false, time.Now()}
	todoMap[id] = todo
	// 렌더링을 사용해서 JSON으로 반환
	rd.JSON(w, http.StatusCreated, todo)
}

func removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)               // mux.Vars를 통해 변수 맵을 만들어 검색
	id, _ := strconv.Atoi(vars["id"]) // id에 있는 문자열을 숫자로 변경
	if _, ok := todoMap[id]; ok {     // todoMap에 id가 있으면
		delete(todoMap, id) // 삭제하고
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                           // mux.Vars를 통해 변수 맵을 만들어 검색
	id, _ := strconv.Atoi(vars["id"])             // id에 있는 문자열을 숫자로 변경
	complete := r.FormValue("complete") == "true" // complete에 true가 담겨오면 체크하라는 뜻이고, false면 체크를 해제하라는 뜻
	if todo, ok := todoMap[id]; ok {
		todo.Completed = complete // true면 체크, false면 해제
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

type Success struct {
	Success bool `json:"success"`
}

func MakeHandler() http.Handler {
	todoMap = make(map[int]*Todo)

	rd = render.New() // 렌더링을 사용해서 JSON 변환을 쉽게하기
	r := mux.NewRouter()

	r.HandleFunc("/todos", getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", removeTodoHandler).Methods("DELETE") // 아이디는 1개 이상인 숫자로만 이루어짐
	r.HandleFunc("/complete-todo/{id:[0-9]+}", completeTodoHandler).Methods("GET")
	// 처음 서버를 시작하고 local:3000으로 들어갔을 때 바로 local:3000/todo.html로 리다이렉트 하는 핸들러
	r.HandleFunc("/", indexHandler)

	return r
}
