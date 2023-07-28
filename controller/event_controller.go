package controller

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/erikrios/my-events/config"
	"github.com/erikrios/my-events/lib/persistence"
	"github.com/gorilla/mux"
)

type SearchCriteria string

const (
	ID   SearchCriteria = "id"
	Name                = "name"
)

func ServeAPI(dbHandler persistence.DatabaseHandler) error {
	handler := NewEventServiceHandler(dbHandler)

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	eventsRouter := r.PathPrefix("/events").Subrouter()

	eventsRouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)

	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)

	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}

type eventServiceHandler struct {
	dbHandler persistence.DatabaseHandler
}

func NewEventServiceHandler(databaseHandler persistence.DatabaseHandler) *eventServiceHandler {
	return &eventServiceHandler{dbHandler: databaseHandler}
}

func (e *eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		writeJson(w, http.StatusBadRequest, &ErrorResponse{Error: "No search criteria found, you can either search by id via /id/4 to search by name via /name/coldplayconcert"})
		return
	}

	searchKey, ok := vars["search"]
	if !ok {
		writeJson(w, http.StatusBadRequest, &ErrorResponse{Error: "error: No search keys found, you can either search by id via /id/4 to search by name via /name/coldplayconcert"})
		return
	}

	var event persistence.Event
	var err error

	switch searchCriteria := SearchCriteria(strings.ToLower(criteria)); searchCriteria {
	case Name:
		event, err = e.dbHandler.FindEventByName(searchKey)
	case ID:
		_, err := hex.DecodeString(searchKey)
		if err != nil {
			writeJson(w, http.StatusBadRequest, &ErrorResponse{Error: "ID Should be in hexadecimal charater"})
			return
		}
		event, err = e.dbHandler.FindEvent([]byte(searchKey))
	default:
		writeJson(w, http.StatusBadRequest, &ErrorResponse{Error: "Unsupported criteria"})
		return
	}

	if err != nil {
		mapErrorResponse(w, err)
		return
	}

	writeJson(w, http.StatusOK, &event)
}

func (e *eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := e.dbHandler.FindAllAvailableEvents()
	if err != nil {
		mapErrorResponse(w, err)
		return
	}

	writeJson(w, http.StatusOK, &events)
}

func (e *eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	var event persistence.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeJson(w, http.StatusBadRequest, &ErrorResponse{Error: "Invalid request body, please check the API documentation"})
		return
	}

	idByte, err := e.dbHandler.AddEvent(event)
	if err != nil {
		mapErrorResponse(w, err)
		return
	}

	idResponse := &IDResponse{ID: string(idByte)}
	writeJson(w, http.StatusCreated, &idResponse)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var internalServerError = &ErrorResponse{
	Error: "Something went wrong.",
}

type IDResponse struct {
	ID string `json:"id"`
}

func writeJson(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(&v)
}

func mapErrorResponse(w http.ResponseWriter, err error) {
	if errors.Is(err, persistence.ErrNotFound) {
		writeJson(w, http.StatusNotFound, &ErrorResponse{Error: "The requested resource not found."})
	} else {
		writeJson(w, http.StatusInternalServerError, internalServerError)
	}
}
