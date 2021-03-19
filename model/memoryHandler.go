package model

import "time"

type memoryHandler struct {
	todoMap map[int]*Todo
}

func newMemoryHandler() DBHandler {
	m := &memoryHandler{}
	m.todoMap = make(map[int]*Todo)
	return m
}

func (m *memoryHandler) Close() {

}

func (m *memoryHandler) GetTodos() []*Todo {
	list := []*Todo{}
	for _, v := range m.todoMap {
		list = append(list, v)
	}
	return list
}
func (m *memoryHandler) AddTodo(name string) *Todo {
	id := len(m.todoMap) + 1 // 임의의 ID
	todo := &Todo{id, name, false, time.Now()}
	m.todoMap[id] = todo
	return todo
}
func (m *memoryHandler) RemoveTodo(id int) bool {
	if _, ok := m.todoMap[id]; ok { // todoMap에 id가 있으면
		delete(m.todoMap, id) // 삭제하고
		return true
	}
	return false
}
func (m *memoryHandler) CompleteTodo(id int, complete bool) bool {
	if todo, ok := m.todoMap[id]; ok {
		todo.Completed = complete // true면 체크, false면 해제
		return true
	}
	return false
}
