package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go-users-service/cmd/plainHTTP/handlers"
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

	userHandlers := handlers.NewUserHandlers(cockroachRepository)

	http.HandleFunc("/user", userHandlers.Handle)

	port := "8081"
	log.Println("Running HTTP in :", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println(err)
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
