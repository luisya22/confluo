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
	Workflows       WorkflowModel
	WorkflowActions WorkFlowActionModel
	Actions         ActionModel
	Providers       ProviderModel
}

func NewModels(db *sqlx.DB) Models {

	return Models{
		Workflows:       WorkflowModel{DB: db},
		WorkflowActions: WorkFlowActionModel{DB: db},
		Actions:         ActionModel{DB: db},
		Providers:       ProviderModel{DB: db},
	}
}
