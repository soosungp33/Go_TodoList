package main

import (
	"log"
	"net/http"

	"github.com/soosungp33/Go_TodoList/app"
	"github.com/urfave/negroni"
)

func main() {
	m := app.MakeHandler()

	n := negroni.Classic() // 미들웨어(기본 파일 서버 기능을 제공 -> public에 있는 파일들을 제공함, 로그 기능도 제공) -> 템플릿 기능을 간단하게
	// mux := pat.New()
	// mux.Handle("/", http.FileServer(http.Dir("public")))과 같은 의미이다.

	n.UseHandler(m) // 래핑

	log.Println("Started App")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
