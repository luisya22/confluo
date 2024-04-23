package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Provider struct {
	Id        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Logo      string    `db:"logo" json:"logo"`
	Actions   []Action  `db:"-" json:"actions"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
	Version   int       `db:"version" json:"version"`
}

type ProviderModel struct {
	DB *sqlx.DB
}

func (model ProviderModel) Insert(p *Provider) error {
	if p.Name == "" {
		return fmt.Errorf("provider name cannot be empty")
	}

	query := `INSERT INTO providers (name, log) VALUES (:name, :logo) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := model.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	return stmt.QueryRowxContext(ctx, p).Scan(&p.Id)
}

func (model ProviderModel) Update(p *Provider) error {
	query := `Update providers SET
		name = :name,
		logo = :logo
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

	err = stmt.QueryRowxContext(ctx, *p).Scan(&p.Version)
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

func (model ProviderModel) Get(id string) (*Provider, error) {
	rowsFound := false

	query := `SELECT p.id, p.name, p.logo, p.version
			a.Id, a.operation 
		FROM providers p
		INNER JOIN actions a ON p.Id = a.Id
		WHERE p.id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := model.DB.QueryxContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var provider Provider

	for rows.Next() {
		rowsFound = true

		var action Action
		err := rows.Scan(
			&provider.Id,
			&provider.Name,
			&provider.Logo,
			&provider.Version,
			&action.Id,
			&action.Operation,
		)
		if err != nil {
			return nil, err
		}

		provider.Actions = append(provider.Actions, action)
	}

	if !rowsFound {
		return nil, ErrRecordNotFound
	}

	return &provider, nil

}
