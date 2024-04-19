package database

import (
	"time"

	"github.com/google/uuid"
)

type Car struct {
	ID        uuid.UUID `db:"id"`
	Mark      string    `db:"mark"`
	Model     string    `db:"model"`
	Color     string    `db:"color"`
	Year      int       `db:"year"`
	RegNums   string    `db:"regNums"`
	Owner     string    `db:"owner"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AddCars struct {
	RegNums []string
}
type CreateCar struct {
	Mark   string
	Model  string
	RegNum string
	Owner  string
	Color  string
	Year   int
}
