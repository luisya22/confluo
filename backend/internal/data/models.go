package data

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflic")
)

type Models struct {
	Workflows WorkflowModel
}

func NewModels(db *sqlx.DB) Models {

	return Models{
		Workflows: WorkflowModel{DB: db},
	}
}
