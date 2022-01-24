package server

import (
	"encoding/json"
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

	value, err := answer.GetAnswerValue(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

func (a *api) postAnswer(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ans answer.Answer
	err = json.Unmarshal(body, &ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = answer.CreateAnswer(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) deleteAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := answer.DeleteAnswer(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) editAnswer(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var ans answer.Answer
	err := json.Unmarshal(body, &ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = answer.EditAnswer(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
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
