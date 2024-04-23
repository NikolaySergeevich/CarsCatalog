package env

import (
	"fmt"
	"log/slog"
	"os"
	"testcar/internal/database/car"
	"testcar/internal/env/config"
)

type Env struct {
	Config         config.Config
	AutoRepository *car.Repository
	Logger         *slog.Logger
}

func InitLogger(cfg *config.Config) (*slog.Logger, error) {
	lvl := new(slog.Level)
	if err := lvl.UnmarshalText([]byte(cfg.Logger.Level)); err != nil {
		return nil, fmt.Errorf("log level UnmarshalText: %w", err)
	}

	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	var handler slog.Handler
	handler = slog.NewTextHandler(os.Stdout, opts)
	if !cfg.Logger.Debug {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler), nil
}