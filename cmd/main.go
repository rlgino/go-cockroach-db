package main

import (
	"context"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	"go-cockroach/internal/handlers"
	"go-cockroach/internal/persistence"
	"log"
	"net/http"
	"os"
)

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
	userHandlers := handlers.NewUserHandlers(cockroachRepository)

	http.HandleFunc("/user", userHandlers.Handle)

	log.Println("Running in :8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
