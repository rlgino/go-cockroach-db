package main

import (
	"context"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
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
	// Add prefix to the router
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv))
	log.Println("Running in port :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
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

func (u userHandlers) DeleteUser(ctx context.Context, params usersvcapi.DeleteUserParams) error {
	err := u.actions.DeleteUser(ctx, params.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (u userHandlers) ListUsers(ctx context.Context) (usersvcapi.Users, error) {
	users, err := u.actions.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	var res = make(usersvcapi.Users, len(users))
	for i, data := range users {
		res[i] = usersvcapi.User{
			ID:   data.ID,
			Name: data.User,
		}
	}
	return res, nil
}

func (u userHandlers) AddUser(ctx context.Context, req *usersvcapi.User) error {
	err := u.actions.CreateUser(ctx, user.Data{
		ID:       req.ID,
		User:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u userHandlers) NewError(_ context.Context, err error) *usersvcapi.ErrorStatusCode {
	return &usersvcapi.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: usersvcapi.Error{
			Message: err.Error(),
		},
	}
}
