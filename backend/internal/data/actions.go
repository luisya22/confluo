package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Action struct {
	Id         string    `db:"id" json:"id"`
	Operation  string    `db:"operation" json:"operation"`
	Provider   Provider  `db:"-" json:"provider"`
	ProviderId string    `db:"provider_id" json:"providerId"`
	CreatedAt  time.Time `db:"created_at" json:"-"`
	UpdatedAt  time.Time `db:"updated_at" json:"-"`
	Version    int       `db:"version" json:"version"`
}

type ActionType int

const (
	ActionTypeTrigger ActionType = iota
	ActionTypeOperation
	ActionTypeConditional
)

type ActionModel struct {
	DB *sqlx.DB
}

func (model ActionModel) Insert(a *Action) error {
	if a.Operation == "" {
		return fmt.Errorf("action operation cannot be empty")
	}

	if a.ProviderId == "" {
		return fmt.Errorf("action provider cannot be empty")
	}

	query := `INSERT INTO actions (operation, provider)
		VALUES (:operation, :provider)
		RETURNIN id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := model.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRowxContext(ctx, a).Scan(&a.Id)
}

func (model ActionModel) Update(a *Action) error {
	query := `UPDATE actions SET
		operation = :operation,
		provider = :provider
		version = version + 1
		WHERE id = :id
		AND version = :version
		RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := model.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, *a).Scan(&a.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
