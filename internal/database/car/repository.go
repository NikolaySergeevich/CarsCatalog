package car

import (
	"context"
	"errors"
	"fmt"
	"testcar/internal/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

func New(auto *pgxpool.Pool, timeout time.Duration) *Repository {
	return &Repository{db: auto, timeout: timeout}
}

type Repository struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req []database.CreateCar) ([]database.Car, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
		INSERT INTO cars (id, mark, model, color, year, regNums, owner, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE
		SET mark = $2, model = $3, color = $5, year = $5, regNums = $6, owner = $7, updated_at = $9
	`
	c := database.Car{}
	res := make([]database.Car, 0, len(req))
	for _, v := range req {
		c.ID = uuid.New()
		c.Mark = v.Mark
		c.Model = v.Model
		c.RegNums = v.RegNum
		c.Color = v.Color
		c.Owner = v.Owner
		c.Year = v.Year

		if _, err := r.db.Exec(ctx, query, uuid.New(), v.Mark, v.Model, v.Color, v.Year, v.RegNum, v.Owner, time.Now(), time.Now()); err != nil {
			var writerErr *pgconn.PgError
			if errors.As(err, &writerErr) && writerErr.Code == "23505" {
				return res, database.ErrConflict
			}
			return res, fmt.Errorf("postgres Exec: %w", err)
		}
		res = append(res, c)
	}
	return res, nil
}