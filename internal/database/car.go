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
	Year      int       `db:"yearCr"`
	RegNums   string    `db:"regNums"`
	Owner     string    `db:"ownerCar"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AddCars struct {
	RegNums []string
}
type CreateCar struct {
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	RegNum string `json:"regNum"`
	Owner  string `json:"owner"`
	Color  string `json:"color"`
	Year   int    `json:"year"`
}
