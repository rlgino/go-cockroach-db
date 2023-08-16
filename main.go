package main

import (
	"context"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"go-cockroach/internal/handlers"
	"go-cockroach/internal/persistence"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func initTable(ctx context.Context, tx pgx.Tx) error {
	log.Println("Drop existing users table if necessary.")
	if _, err := tx.Exec(ctx, "DROP TABLE IF EXISTS users"); err != nil {
		return err
	}

	log.Println("Creating users table.")
	if _, err := tx.Exec(ctx,
		"CREATE TABLE users (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), userName varchar, password varchar)"); err != nil {
		return err
	}
	return nil
}

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
		return initTable(context.Background(), tx)
	})
	cockroachRepository := persistence.New(conn)
	userHandlers := handlers.NewUserHandlers(cockroachRepository)

	http.HandleFunc("/user", userHandlers.Handle)

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
