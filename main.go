package main

import (
	"log"
	"net/http"

	"github.com/alexrodfe/golang-api/answer"
	"github.com/alexrodfe/golang-api/server"
)

func main() {
	s := server.New()

	answer.InitAnswers()
	answer.InitEvents()

	log.Fatal(http.ListenAndServe(":8081", s.Router()))
}
