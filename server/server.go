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
	anse   answer.AnswerEngine
}

type Server interface {
	Router() http.Handler
}

func New(anse answer.AnswerEngine) Server {
	a := &api{anse: anse}

	r := mux.NewRouter()
	r.HandleFunc("/answer", a.postAnswer).Methods(http.MethodPost)
	r.HandleFunc("/answer/{key:[a-zA-Z1-9_]+}", a.getAnswer).Methods(http.MethodGet)
	r.HandleFunc("/answer/{key:[a-zA-Z1-9_]+}", a.deleteAnswer).Methods(http.MethodDelete)
	r.HandleFunc("/answer/{key:[a-zA-Z1-9_]+}", a.editAnswer).Methods(http.MethodPost)
	r.HandleFunc("/answer/history/{key:[a-zA-Z1-9_]+}", a.getAnswerHistory).Methods(http.MethodGet)

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

func (a *api) getAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := a.anse.GetAnswerValue(key)
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

	err = a.anse.CreateAnswer(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) deleteAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := a.anse.DeleteAnswer(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) editAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

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

	if key != ans.Key {
		http.Error(w, "answer key does not match requested path key", http.StatusBadRequest)
		return
	}

	err = a.anse.EditAnswer(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *api) getAnswerHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	events, err := a.anse.GetAnswerHistory(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
