package persistence

import (
	"context"
	"log"
	"time"

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
	return repo
}

// InitRepository will init the repository and create every part of the struct
func InitRepository(ctx context.Context, tx pgx.Tx) error {
	// Dropping existing table if it exists
	log.Println("Drop existing users table if necessary.")
	if _, err := tx.Exec(ctx, "DROP TABLE IF EXISTS users"); err != nil {
		return err
	}

	// Create the accounts table
	log.Println("Creating users table.")
	if _, err := tx.Exec(ctx,
		"CREATE TABLE users (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), firstName varchar, lastName varchar, birthdate varchar, status varchar)"); err != nil {
		return err
	}
	return nil
}

const dateFormat = "2006-01-02"

func (repo *CockroachRepository) SaveUser(ctx context.Context, user user.Data) error {
	log.Println("Creating new row...")
	if _, err := repo.conn.Exec(ctx,
		"INSERT INTO users (id, firstName, lastName, birthdate, status) VALUES ($1, $2, $3, $4, $5)", user.ID, user.FistName, user.LastName, user.Birthdate.Format(dateFormat), user.Status); err != nil {
		return err
	}
	return nil
}

func (repo *CockroachRepository) FindUserByID(ctx context.Context, id uuid.UUID) (user.Data, error) {
	row := repo.conn.QueryRow(ctx, "SELECT id, firstName, lastName, birthdate, status FROM users where id = $1", id.String())
	var ret user.Data
	var birthdate string
	if err := row.Scan(&ret.ID, &ret.FistName, &ret.LastName, &birthdate, &ret.Status); err != nil {
		return user.Data{}, err
	}
	formattedBirthdate, _ := time.Parse(dateFormat, birthdate)
	ret.Birthdate = formattedBirthdate
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
