package persistence

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	User     string    `json:"user,omitempty"`
	Password string    `json:"password,omitempty"`
}

type CockroachRepository struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *CockroachRepository {
	return &CockroachRepository{
		conn,
	}
}

func (repo *CockroachRepository) SaveUser(ctx context.Context, user User) error {
	log.Println("Creating new row...")
	if _, err := repo.conn.Exec(ctx,
		"INSERT INTO users (id, userName, password) VALUES ($1, $2, $3)", user.ID, user.User, user.Password); err != nil {
		return err
	}
	return nil
}

func (repo *CockroachRepository) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := repo.conn.Query(ctx, "SELECT id, userName FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []User
	for rows.Next() {
		var id uuid.UUID
		var userName string
		if err := rows.Scan(&id, &userName); err != nil {
			return nil, err
		}
		log.Printf("%s: %s\n", id, userName)
		ret = append(ret, User{
			ID:   id,
			User: userName,
		})
	}
	return ret, nil
}

func (repo *CockroachRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Delete two rows into the "accounts" table.
	log.Printf("Deleting rows with ID %s", id)
	if _, err := repo.conn.Exec(ctx,
		"DELETE FROM users WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}
