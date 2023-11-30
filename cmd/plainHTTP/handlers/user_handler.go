package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-users-service/internal/core/user"
	"log"
	"net/http"
)

type UserHandlers struct {
	actions user.Actions
}

func NewUserHandlers(repo user.Repository) UserHandlers {
	return UserHandlers{
		user.NewActions(repo),
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
	var userToDelete user.Data
	err := json.NewDecoder(req.Body).Decode(&userToDelete)
	if err != nil {
		writeInternalError(res, "Invalid request", err)
		return
	}
	err = handler.actions.DeleteUser(req.Context(), userToDelete.ID)
	if err != nil {
		writeInternalError(res, "error deleting user", err)
		return
	}
	writeSuccess(res, nil)
}

func (handler *UserHandlers) createUser(res http.ResponseWriter, req *http.Request) {
	var userToCreate user.Data
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
	err = handler.actions.CreateUser(req.Context(), userToCreate)
	if err != nil {
		writeInternalError(res, "error creating user", err)
		return
	}
	writeSuccess(res, nil)

}

func (handler *UserHandlers) listUsers(res http.ResponseWriter, req *http.Request) {
	users, err := handler.actions.ListUsers(req.Context())
	if err != nil {
		writeInternalError(res, "error listing user", err)
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
	res.WriteHeader(http.StatusInternalServerError)
	errMsg := fmt.Sprintf(msg+": %s", err.Error())
	log.Println(errMsg)
	_, errWriting := res.Write([]byte(errMsg))
	if errWriting != nil {
		log.Println(errWriting)
	}
}
