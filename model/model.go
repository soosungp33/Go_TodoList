package model

import "time"

// db부분 구현하는 패키지
type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type DBHandler interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(sessionId string, name string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, complete bool) bool
	Close()
}

func NewDBHandler(dbConn string) DBHandler {
	// return newMemoryHandler() -> 인메모리 db
	// return newSqliteHandler(filepath) -> sqlite
	return newPQHandler(dbConn) // -> postgreDB
}
