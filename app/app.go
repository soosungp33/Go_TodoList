package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MakeHandler() http.Handler {
	r := mux.NewRouter()

	return r
}
