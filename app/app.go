package app

import (
	"net/http"
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

func adddTodoHandler(w http.ResponseWriter, r *http.Request) { // 프론트에서 올 때 name에다 Item을 넣어서 온다.
	name := r.FormValue("name") // 키를 통해서 value에 있는 Item을 꺼낸다.
	id := len(todoMap) + 1      // 임의의 ID
	todo := &Todo{id, name, false, time.Now()}
	todoMap[id] = todo
	// 렌더링을 사용해서 JSON으로 반환
	rd.JSON(w, http.StatusOK, todo)
}

func addTestTodos() { // 테스트용 데이터
	todoMap[1] = &Todo{1, "Test1", false, time.Now()}
	todoMap[2] = &Todo{2, "Test2", true, time.Now()}
	todoMap[3] = &Todo{3, "Test3", false, time.Now()}
}

func MakeHandler() http.Handler {
	todoMap = make(map[int]*Todo)
	addTestTodos()

	rd = render.New() // 렌더링을 사용해서 JSON 변환을 쉽게하기
	r := mux.NewRouter()

	r.HandleFunc("/todos", getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", adddTodoHandler).Methods("POST")
	r.HandleFunc("/", indexHandler) // 처음 서버를 시작하고 local:3000으로 들어갔을 때 바로 local:3000/todo.html로 리다이렉트 하는 핸들러

	return r
}
