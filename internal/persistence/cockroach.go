package persistence

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go-users-service/internal/core/user"
)

type CockroachRepository struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *CockroachRepository {
	repo := &CockroachRepository{
		conn,
	}
	id := uuid.New()
	log.Println("Using ID ", id)
	repo.SaveUser(context.Background(), user.Data{
		ID:       id,
		User:     "test1",
		Password: "prueba",
	})
	return repo
}

// InitRepository will init the repository and create every part of the struct
func InitRepository(ctx context.Context, tx pgx.Tx) error {
	// TODO: Adding creating data
	// Dropping existing table if it exists
	log.Println("Drop existing users table if necessary.")
	if _, err := tx.Exec(ctx, "DROP TABLE IF EXISTS users"); err != nil {
		return err
	}

	// Create the accounts table
	log.Println("Creating users table.")
	if _, err := tx.Exec(ctx,
		"CREATE TABLE users (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), userName varchar, password varchar)"); err != nil {
		return err
	}
	return nil
}

func (repo *CockroachRepository) SaveUser(ctx context.Context, user user.Data) error {
	log.Println("Creating new row...")
	if _, err := repo.conn.Exec(ctx,
		"INSERT INTO users (id, userName, password) VALUES ($1, $2, $3)", user.ID, user.User, user.Password); err != nil {
		return err
	}
	return nil
}

func (repo *CockroachRepository) ListUsers(ctx context.Context) ([]user.Data, error) {
	rows, err := repo.conn.Query(ctx, "SELECT id, userName FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []user.Data
	for rows.Next() {
		var id uuid.UUID
		var userName string
		if err := rows.Scan(&id, &userName); err != nil {
			return nil, err
		}
		ret = append(ret, user.Data{
			ID:   id,
			User: userName,
		})
	}
	return ret, nil
}

func (repo *CockroachRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Delete two rows into the "accounts" table.
	if _, err := repo.conn.Exec(ctx,
		"DELETE FROM users WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}
