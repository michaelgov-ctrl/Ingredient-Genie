package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), logLevel(app.config.logLevel)),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info("starting server", "addr", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)

	return nil
}

func logLevel(s string) slog.Level {
	m := map[string]slog.Level{
		"trace":   slog.Level(-8),
		"debug":   slog.LevelDebug,
		"info":    slog.LevelInfo,
		"warning": slog.LevelWarn,
		"error":   slog.LevelError,
	}

	if level, ok := m[s]; ok {
		return level
	}

	return slog.LevelError
}

//go:embed meals.sqlite
var embeddedDatabaseBytes []byte

func hydrateDB(cfg config) error {
	path, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	dbPath := filepath.Join(dir, cfg.db.path)

	_, err = os.Stat(dbPath)
	if err == nil {
		return nil
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	return os.WriteFile(dbPath, embeddedDatabaseBytes, 0600)
}

func openDB(cfg config) (*sql.DB, error) {
	if err := hydrateDB(cfg); err != nil {
		return nil, err
	}

	// TODO: check if path is relative or absolute
	// if relative and is an end node (file) with no
	// preceding directory, maybase assume it's in
	// the same directory as this executable
	db, err := sql.Open("sqlite", cfg.db.path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
