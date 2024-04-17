package data

import "time"

type WorkflowAction struct {
	Id           string                 `db:"id" json:"id"`
	Text         string                 `db:"text" json:"text"`
	WorkflowId   string                 `db:"workflow_id" json:"workflowId"`
	ActionId     string                 `db:"action_id" json:"actionId"`
	Action       Action                 `db:"-" json:"action"`
	Type         string                 `db:"type" json:"type"`
	NextActionId string                 `db:"next_action_id" json:"next_action_id"`
	NextAction   string                 `db:"-" json:"nextAction"`
	Params       map[string]interface{} `db:"params" json:"params"`
	CreatedAt    time.Time              `db:"created_at" json:"-"`
	UpdatedAt    time.Time              `db:"updated_at" json:"-"`
	Version      int                    `db:"version" json:"version"`
}
