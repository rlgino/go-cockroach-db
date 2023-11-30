package ogeserver

import (
	"go-users-service/cmd/ogen/usersvcapi"
	"go-users-service/internal/core/user"
	"log"
	"net/http"
)

func NewServer(repo user.Repository) *OgenServer {
	// Create service instance.
	service := NewHandlers(repo)
	// Create generated server.
	srv, err := usersvcapi.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	// Add prefix to the router
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv))
	return &OgenServer{
		mux,
	}
}

type OgenServer struct {
	mux *http.ServeMux
}

func (srv *OgenServer) ListenAndServe(port string) {
	log.Println("Running in port ", port)
	if err := http.ListenAndServe(":"+port, srv.mux); err != nil {
		log.Fatal(err)
	}
}
