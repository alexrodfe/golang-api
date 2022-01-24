package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexrodfe/golang-api/answer"
	"github.com/gorilla/mux"
)

type api struct {
	router http.Handler
}

type Server interface {
	Router() http.Handler
}

func (a *api) getAnswers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(answer.AllAnswersIndexed)
}

func (a *api) getAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := answer.GetValue(key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

func (a *api) postAnswer(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // TODO : err handling

	var ans answer.Answer
	err := json.Unmarshal(body, &ans)
	if err != nil {
		fmt.Println("ouch")
	}

	answer.PostValue(ans)
	w.WriteHeader(http.StatusOK)
}

func (a *api) deleteAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	answer.DeleteValue(key)

	w.WriteHeader(http.StatusOK)
}

func (a *api) editAnswer(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // TODO : err handling

	var ans answer.Answer
	err := json.Unmarshal(body, &ans)
	if err != nil {
		fmt.Println("ouch")
	}

	err = answer.EditValue(ans)
	if err != nil {
		fmt.Println("ouch")
	}
	w.WriteHeader(http.StatusOK)
}

func New() Server {
	a := &api{}

	r := mux.NewRouter()
	r.HandleFunc("/answers", a.getAnswers).Methods(http.MethodGet)
	r.HandleFunc("/answer/{key:[a-zA-Z_]+}", a.getAnswer).Methods(http.MethodGet)
	r.HandleFunc("/answer", a.postAnswer).Methods(http.MethodPost)
	r.HandleFunc("/answer/{key:[a-zA-Z_]+}", a.deleteAnswer).Methods(http.MethodDelete)
	r.HandleFunc("/answer/{key:[a-zA-Z_]+}", a.editAnswer).Methods(http.MethodPost)

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}
