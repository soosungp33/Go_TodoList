package main

import (
	"log"
	"net/http"
	"os"

	"github.com/soosungp33/Go_TodoList/app"
)

func main() {
	port := os.Getenv("PORT") // heroku에 배포하기 위한 port
	// m := app.MakeHandler("./test.db") // 나중에 실행인자로 파일명을 넣어줄 수 있음 -> sqlite db 사용
	m := app.MakeHandler(os.Getenv("DATABASE_URL")) // postgre db 사용(heroku에서 발급받은 postgreDB URL을 설정)
	defer m.Close()                                 // 앱이 종료되기 전에 db를 종료해준다.(main에서 컨트롤해야한다.)

	log.Println("Started App")
	err := http.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}
}
