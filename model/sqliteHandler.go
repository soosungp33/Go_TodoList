package model

import (
	"database/sql"
	"time"

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
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() { // 다음 행으로 가면서 레코드를 읽음(모든 데이터)
		var todo Todo
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt) // 레코드를 Todo구조체에 넣어준다.
		todos = append(todos, &todo)
	}
	return todos
}
func (s *sqliteHandler) AddTodo(name string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (name, completed, createdAt) VALUES (?, ?, datetime('now'))") // 쿼리문 작성
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(name, false) // 쿼리문 실행(? 아규먼트 순서대로)
	if err != nil {
		panic(err)
	}
	id, _ := rst.LastInsertId() // 마지막으로 추가된 레코드의 id를 알려준다.
	var todo Todo
	todo.ID = int(id)
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()

	return &todo
}
func (s *sqliteHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?") // 쿼리문 작성
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(id) // 쿼리문 실행(? 아규먼트 순서대로)
	if err != nil {
		panic(err)
	}
	cnt, err := rst.RowsAffected() // 쿼리문으로 영향받은 레코드 갯수를 반환
	if err != nil {
		panic(err)
	}
	return cnt > 0 // 변경된 내용이 1개면 true, 0개면 false
}
func (s *sqliteHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?") // 쿼리문 작성
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(complete, id) // 쿼리문 실행(? 아규먼트 순서대로)
	if err != nil {
		panic(err)
	}
	cnt, err := rst.RowsAffected() // 쿼리문으로 영향받은 레코드 갯수를 반환
	if err != nil {
		panic(err)
	}

	return cnt > 0 // 변경된 내용이 1개면 true, 0개면 false
}
