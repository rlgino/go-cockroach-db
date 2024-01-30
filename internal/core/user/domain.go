package user

import (
	"github.com/google/uuid"
	"time"
)

type Data struct {
	ID        uuid.UUID `json:"id"`
	FistName  string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Birthdate time.Time `json:"birthdate"`
	Status    string    `json:"status"`
}
