package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/luisya22/confluo/backend/internal/data"
	"github.com/luisya22/confluo/backend/internal/executor"
	"github.com/luisya22/confluo/backend/oauth"
)

type Application struct {
	config       Config
	logger       *slog.Logger
	models       data.Models
	wg           sync.WaitGroup
	executor     *executor.Executor
	oauhtService *oauth.OauthService
}

func NewApplication(cfg Config) *Application {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDb(cfg)
	if err != nil {
		log.Fatal(err)
	}

	models := data.NewModels(db)

	oauthConfig := oauth.Config{
		Github: oauth.Github{
			ClientId:     cfg.Providers.Github.ClientId,
			ClientSecret: cfg.Providers.Github.ClientSecret,
			Url:          cfg.Providers.Github.Url,
			RedirectUrl:  cfg.Providers.Github.RedirectUrl,
			UserUrl:      cfg.Providers.Github.UserUrl,
		},
	}

	oauthService := oauth.NewOauthService(oauthConfig)

	return &Application{
		config:       cfg,
		logger:       logger,
		models:       models,
		wg:           sync.WaitGroup{},
		oauhtService: oauthService,
	}

}

func openDb(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Application) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks", "addr", srv.Addr)

		app.wg.Wait()

		shutdownError <- nil
	}()

	app.logger.Info(
		"starting server",
		"addr", srv.Addr,
		"env", app.config.Env,
	)

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
