package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // 암시적으로 사용
)

type sqliteHandler struct {
	db *sql.DB // sql 인스턴스를 가지고 있어야 한다.
}

func newSqliteHandler() DBHandler {
	db, err := sql.Open("sqlite3", "./test.db") // test.db라는 파일db를 연다.
	if err != nil {
		panic(err)
	}
	statement, _ := db.Prepare( // 쿼리문 작성
		`CREATE TABLE IF NOT EXISTS todos (
			id 		  INTEGER PRIMARY KEY AUTOINCREMENT,
			name 	  TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		)`)
	statement.Exec() // 퀴리문을 실행

	return &sqliteHandler{db} // 이 db를 계속 사용해야하니까 인스턴스로 저장하고 리턴해준다.
}

func (s *sqliteHandler) Close() { // db가 사라지기 전에 닫아줘야해서 만드는 함수
	s.db.Close()
}

func (s *sqliteHandler) GetTodos() []*Todo {
	return nil
}
func (s *sqliteHandler) AddTodo(name string) *Todo {
	return nil
}
func (s *sqliteHandler) RemoveTodo(id int) bool {
	return false
}
func (s *sqliteHandler) CompleteTodo(id int, complete bool) bool {
	return false
}
