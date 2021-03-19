package main

import (
	"log"
	"net/http"

	"github.com/soosungp33/Go_TodoList/app"
	"github.com/urfave/negroni"
)

func main() {
	m := app.MakeHandler()
	defer m.Close() // 앱이 종료되기 전에 db를 종료해준다.
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("Started App")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
