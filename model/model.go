package model

import "time"

// db부분 구현하는 패키지
type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var todoMap map[int]*Todo // 인메모리 db역할

func init() { // 이 패키지가 initialize될 때 한 번만 호출됨(초기화)
	todoMap = make(map[int]*Todo)
}

func GetTodos() []*Todo {
	list := []*Todo{}
	for _, v := range todoMap {
		list = append(list, v)
	}
	return list
}

func AddTodo(name string) *Todo {
	id := len(todoMap) + 1 // 임의의 ID
	todo := &Todo{id, name, false, time.Now()}
	todoMap[id] = todo
	return todo
}

func RemoveTodo(id int) bool {
	if _, ok := todoMap[id]; ok { // todoMap에 id가 있으면
		delete(todoMap, id) // 삭제하고
		return true
	}
	return false
}

func CompleteTodo(id int, complete bool) bool {
	if todo, ok := todoMap[id]; ok {
		todo.Completed = complete // true면 체크, false면 해제
		return true
	}
	return false
}
