package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go-users-service/cmd/ogen/usersvcapi"
	"go-users-service/internal/core/user"
	"go-users-service/internal/persistence"
	"log"
	"net/http"
	"os"
)

func main() {
	conn, cockroachRepository, err := CreateConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	// Create service instance.
	service := NewHandlers(cockroachRepository)
	// Create generated server.
	srv, err := usersvcapi.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	// Add prefix to the router
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv))

	port := "8080"
	log.Println("Running in port ", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func CreateConnection() (*pgx.Conn, user.Repository, error) {
	// Read in connection string
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, nil, err
	}
	config.RuntimeParams["application_name"] = "$ docs_simplecrud_gopgx"
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, nil, err
	}
	return conn, persistence.New(conn), nil
}
