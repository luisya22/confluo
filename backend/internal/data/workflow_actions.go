package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type WorkflowAction struct {
	Id           string                 `db:"id" json:"id"`
	Text         string                 `db:"text" json:"text"`
	WorkflowId   string                 `db:"workflow_id" json:"workflowId"`
	ActionId     string                 `db:"action_id" json:"actionId"`
	Action       Action                 `db:"-" json:"action"`
	Type         string                 `db:"type" json:"type"`
	NextActionId sql.NullString         `db:"next_action_id" json:"next_action_id"`
	NextAction   string                 `db:"-" json:"nextAction"`
	Params       map[string]interface{} `db:"params" json:"params"`
	CreatedAt    time.Time              `db:"created_at" json:"-"`
	UpdatedAt    time.Time              `db:"updated_at" json:"-"`
	Version      int                    `db:"version" json:"version"`
}

type WorkFlowActionModel struct {
	DB *sqlx.DB
}

func (model WorkFlowActionModel) Insert(wa *WorkflowAction) error {
	if wa.WorkflowId == "" {
		return fmt.Errorf("workflow id cannot be empty")
	}

	if wa.ActionId == "" {
		return fmt.Errorf("action id cannot be empty")
	}

	query := `INSERT INTO workflow_actions (text, workflow_id, action_id, type)
		VALUES (:text, :workflow_id, :action_id, :type)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := model.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRowxContext(ctx, wa).Scan(&wa.Id)
}

func (model WorkFlowActionModel) Update(wa *WorkflowAction) error {

	paramsJSON, err := json.Marshal(wa.Params)
	if err != nil {
		return err
	}

	query := `UPDATE workflow_actions SET
		text = :text,
		params = :params,
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

	paramMap := map[string]interface{}{
		"id":      wa.Id,
		"text":    wa.Text,
		"params":  paramsJSON,
		"version": wa.Version,
	}

	err = stmt.QueryRowxContext(ctx, paramMap).Scan(&wa.Version)
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
