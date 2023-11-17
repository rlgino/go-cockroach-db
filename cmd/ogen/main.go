package main

import (
	"context"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go-cockroach/cmd/ogen/usersvcapi"
	"go-cockroach/internal/core/user"
	"go-cockroach/internal/persistence"
	"log"
	"net/http"
	"os"
)

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest -package usersvcapi --target usersvcapi --clean usersvc-oas.yml

func main() {
	// Read in connection string
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	config.RuntimeParams["application_name"] = "$ docs_simplecrud_gopgx"
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return persistence.InitRepository(context.Background(), tx)
	})
	if err != nil {
		log.Fatal(err)
	}
	cockroachRepository := persistence.New(conn)

	// Create service instance.
	service := NewHandlers(cockroachRepository)
	// Create generated server.
	srv, err := usersvcapi.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}

func NewHandlers(repo user.Repository) usersvcapi.Handler {
	return userHandlers{
		actions: user.NewActions(repo),
	}
}

type userHandlers struct {
	actions user.Actions
}

func (u userHandlers) AddUser(ctx context.Context, req *usersvcapi.User) error {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return err
	}
	err = u.actions.CreateUser(ctx, user.Data{
		ID:       id,
		User:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u userHandlers) NewError(ctx context.Context, err error) *usersvcapi.ErrorStatusCode {
	return &usersvcapi.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: usersvcapi.Error{
			Code:    123,
			Message: err.Error(),
		},
	}
}
