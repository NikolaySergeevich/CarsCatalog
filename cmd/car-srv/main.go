package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"testcar/internal/database/car"
	"testcar/internal/env/config"
	"testcar/internal/handler"
	"testcar/pkg/api/carapi"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sethvargo/go-envconfig"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	if er := migration(); er != nil {
		log.Fatal(er)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := runMain(ctx); err != nil {
		log.Fatal(err)
	}
}

func runMain(ctx context.Context) error {
	router := chi.NewRouter()

	var cfg config.PostgresConfig
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return fmt.Errorf("env processing: %w", err)
	}

	carDBConn, err := pgxpool.Connect(ctx, cfg.ConnectionURL())
	if err != nil {
		return fmt.Errorf("pgxpool Connect: %w", err)
	}
	log.Println("connect - " + cfg.ConnectionURL())

	carRepository := car.New(carDBConn, 5*time.Second)
	handler := handler.NewHandler(*carRepository)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		router.Mount(
			"/api", carapi.HandlerWithOptions(
				handler, carapi.ChiServerOptions{
					BaseURL: "/v1",
					ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
						slog.Error("handle error", slog.String("err", err.Error()))
						return
					},
				},
			),
		)

		srv := http.Server{
			Addr:              ":8111",
			Handler:           router,
			ReadTimeout:       20 * time.Second,
			ReadHeaderTimeout: 20 * time.Second,
			WriteTimeout:      20 * time.Second,
			IdleTimeout:       20 * time.Second,
			MaxHeaderBytes:    10 * 1024 * 1024, // 10mib
		}
		go func() {
			<-ctx.Done()
			// если посылаем сигнал завершения то завершаем работу нашего сервера
			srv.Close()
		}()
		slog.Info(fmt.Sprintf("http server was started %s", ":8111"))
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("http.Server ListenAndServe", slog.String("err", err.Error()))
			return
		}
	}()

	wg.Wait()
	return nil

}

func migration() error {
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5434/auto")
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	m, err := migrate.New("file://../../migrations", "postgres://localhost:5434/auto?sslmode=disable&user=postgres&password=postgres")
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("m.Up: %w", err)
	}
	log.Println("Migration completed")
	return nil
}
