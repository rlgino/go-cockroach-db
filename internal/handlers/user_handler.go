package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-cockroach/internal/persistence"
	"log"
	"net/http"
)

type UserHandlers struct {
	repo *persistence.CockroachRepository
}

func NewUserHandlers(repo *persistence.CockroachRepository) UserHandlers {
	return UserHandlers{
		repo: repo,
	}
}

func (handler *UserHandlers) Handle(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		handler.listUsers(res, req)
	case http.MethodPost:
		handler.createUser(res, req)
	case http.MethodDelete:
		handler.deleteUser(res, req)
	default:
		writeBadRequest(res, "Method not handled")
	}
}

func (handler *UserHandlers) deleteUser(res http.ResponseWriter, req *http.Request) {
	var userToDelete persistence.User
	err := json.NewDecoder(req.Body).Decode(&userToDelete)
	if err != nil {
		writeInternalError(res, "Invalid request", err)
		return
	}
	err = handler.repo.DeleteUser(req.Context(), userToDelete.ID)
	if err != nil {
		writeInternalError(res, "error deleting user", err)
		return
	}
	writeSuccess(res, nil)
}

func (handler *UserHandlers) createUser(res http.ResponseWriter, req *http.Request) {
	var userToCreate persistence.User
	err := json.NewDecoder(req.Body).Decode(&userToCreate)
	if err != nil {
		writeInternalError(res, "Invalid request", err)
		return
	}
	if userToCreate.User == "" {
		writeBadRequest(res, "User required")
		return
	}
	if userToCreate.Password == "" {
		writeBadRequest(res, "Password required")
		return
	}
	userToCreate.ID, err = uuid.NewRandom()
	if err != nil {
		writeInternalError(res, "Error creating UUID", err)
		return
	}
	err = handler.repo.SaveUser(req.Context(), userToCreate)
	if err != nil {
		writeInternalError(res, "error creating user", err)
		return
	}
	writeSuccess(res, nil)

}

func (handler *UserHandlers) listUsers(res http.ResponseWriter, req *http.Request) {
	users, err := handler.repo.ListUsers(req.Context())
	if err != nil {
		writeInternalError(res, "error listing users", err)
		return
	}
	writeSuccess(res, users)

}

func writeSuccess(res http.ResponseWriter, body interface{}) {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	if body == nil {
		return
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		writeInternalError(res, "Error unmarshalling", err)
		return
	}
	_, err = res.Write(bytes)
	if err != nil {
		log.Println(err)
	}
}

func writeBadRequest(res http.ResponseWriter, msg string) {
	_, err := res.Write([]byte(msg))
	if err != nil {
		log.Println(err)
	}
	res.WriteHeader(http.StatusNotFound)
}

func writeInternalError(res http.ResponseWriter, msg string, err error) {
	_, errWriting := res.Write([]byte(fmt.Sprintf(msg+": %s", err)))
	if errWriting != nil {
		log.Println(errWriting)
	}
	res.WriteHeader(http.StatusInternalServerError)
}
