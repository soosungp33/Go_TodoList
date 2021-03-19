package main

import (
	"log"
	"net/http"

	"github.com/soosungp33/Go_TodoList/app"
	"github.com/urfave/negroni"
)

func main() {
	m := app.MakeHandler("./test.db") // 나중에 실행인자로 파일명을 넣어줄 수 있음
	defer m.Close()                   // 앱이 종료되기 전에 db를 종료해준다.(main에서 컨트롤해야한다.)
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("Started App")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
