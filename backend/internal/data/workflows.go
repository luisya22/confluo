package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Workflow struct {
	Id        string           `db:"id" json:"id"`
	UserId    string           `db:"user_id" json:"user_id"`
	Name      string           `db:"name" json:"name"`
	TriggerId sql.NullString   `db:"trigger_id" json:"trigger_id"`
	Trigger   WorkflowAction   `db:"-" json:"trigger"`
	Actions   []WorkflowAction `db:"-" json:"actions"`
	CreatedAt time.Time        `db:"created_at" json:"-"`
	UpdatedAt time.Time        `db:"updated_at" json:"-"`
	Version   int              `db:"version" json:"version"`
}

type WorkflowModel struct {
	DB *sqlx.DB
}

func (wm WorkflowModel) Insert(w *Workflow) error {
	if w.UserId == "" {
		return fmt.Errorf("user id cannot be empty")
	}

	if w.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	query := `INSERT INTO workflows (name, user_id) values (:name, :user_id) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := wm.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	return stmt.QueryRowxContext(ctx, w).Scan(&w.Id)
}

func (wm WorkflowModel) Update(w *Workflow) error {
	query := `UPDATE workflows SET 
			name = :name,
			version = version + 1
		WHERE id = :id
		AND version = :version
		RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := wm.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, *w).Scan(&w.Version)
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

func (wm WorkflowModel) Get(id string) (*Workflow, error) {
	// TODO: Join Action
	query := `SELECT 
			workflows.id, workflows.name, workflows.trigger_id, workflows.version,
			workflow_actions.id, workflow_actions.text, workflow_actions.type, workflow_actions.next_action_id, workflow_actions.params 
			actions.id, actions.provider, actions.operation
		FROM workflows
		LEFT JOIN workflow_actions on workflows.id = workflow_actions.workflow_id
		LEFT JOIN actions on workflow_actions.action_id = actions.id
		WHERE workflows.id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := wm.DB.QueryxContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	var workflow Workflow
	actions := []WorkflowAction{}

	for rows.Next() {

		var workflowAction WorkflowAction
		rows.Scan(
			&workflow.Id,
			&workflow.Name,
			&workflow.TriggerId,
			&workflow.Version,
			&workflowAction.Text,
			&workflowAction.Type,
			&workflowAction.NextActionId,
			&workflowAction.Params,
			&workflowAction.Action.Id,
			&workflowAction.Action.Provider,
			&workflowAction.Action.Operation,
		)

		actions = append(actions, workflowAction)
	}

	workflow.Actions = actions

	return &workflow, nil
}

func (wm WorkflowModel) GetAll(userId string, filters Filters) ([]*Workflow, Metadata, error) {
	query := `SELECT w.id, w.name. w.trigger_id, w.version, STRING_AGG(a.provider, ',')
		FROM workflows w
		LEFT JOIN workflow_actions wa ON w.id = wa.workflow_id
		LEFT JOIN actions a ON wa.action_id = a.id
		WHERE uerId = $1
		GROUP BY w.id	
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := wm.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	workflows := []*Workflow{}

	for rows.Next() {
		var workflow Workflow
		var actionsStr string
		var actions []WorkflowAction

		err := rows.Scan(
			&workflow.Id,
			&workflow.Name,
			&workflow.TriggerId,
			&workflow.Version,
			&actionsStr,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		actionsSpl := strings.Split(actionsStr, ",")

		for _, action := range actionsSpl {
			a := WorkflowAction{
				Action: Action{
					Provider: action,
				},
			}
			actions = append(actions, a)
		}

		workflow.Actions = actions

		workflows = append(workflows, &workflow)

	}

	return workflows, Metadata{}, nil
}

func (wm WorkflowModel) Delete(id string) error {
	query := `DELETE FROM workflows WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := wm.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
