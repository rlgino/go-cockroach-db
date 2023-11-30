package main

import (
	"context"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	"go-users-service/cmd/grpcserver"
	ogeserver "go-users-service/cmd/ogen"
	"go-users-service/internal/persistence"
	"log"
	"os"
	"sync"
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

	wg := sync.WaitGroup{}
	// TODO: Add failed case with a channel
	wg.Add(1)
	go func() {
		// Start GRPC
		var grpcServer Server
		grpcServer = grpcserver.New(cockroachRepository)
		grpcServer.ListenAndServe("8080")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		var httpServer Server
		httpServer = ogeserver.NewServer(cockroachRepository)
		httpServer.ListenAndServe("8080")
	}()

	wg.Wait()
}

type Server interface {
	ListenAndServe(port string)
}
