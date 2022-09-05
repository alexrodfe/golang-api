package main

import (
	"log"
	"net/http"

	"github.com/alexrodfe/golang-api/api/answer"
	"github.com/alexrodfe/golang-api/api/server"
)

var AllAnswersIndexed answer.MapOfAnswers
var AllEventsIndexed answer.MapOfEvents

func main() {
	AllAnswersIndexed = answer.InitAnswers()
	AllEventsIndexed = answer.InitEvents()

	anse := answer.NewAnswerEngine(&AllAnswersIndexed, &AllEventsIndexed)
	s := server.New(anse)

	log.Fatal(http.ListenAndServe(":8080", s.Router()))
}
