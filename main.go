package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
}

func ServeAPI(endpoint string) error {
	handler := &eventServiceHandler{}
	r := mux.NewRouter()
	eventsRouter := r.PathPrefix("/events").Subrouter()

	eventsRouter.Methods(http.MethodGet).Path("/{search_criteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsRouter.Methods(http.MethodGet).Path("").HandlerFunc(handler.allEventHandler)
	eventsRouter.Methods(http.MethodPost).Path("").HandlerFunc(handler.newEventHandler)

	return http.ListenAndServe(endpoint, r)
}

type eventServiceHandler struct{}

func (eh *eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {

}

func (eh *eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {

}

func (eh *eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {

}
