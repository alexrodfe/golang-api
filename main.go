package main

import (
	"log"
	"net/http"

	"github.com/alexrodfe/golang-api/answer"
	"github.com/alexrodfe/golang-api/server"
)

func main() {
	s := server.New()

	answer.AllAnswersIndexed = make(answer.MapOfAnswers)
	answer.AllAnswersIndexed["key"] = "hola"

	log.Fatal(http.ListenAndServe(":8081", s.Router()))
}
