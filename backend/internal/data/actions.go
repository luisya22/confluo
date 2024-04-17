package data

import "time"

type Action struct {
	Id        string    `db:"id" json:"id"`
	Operation string    `db:"operation" json:"operation"`
	Provider  string    `db:"provider" json:"provider"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
	Version   int       `db:"version" json:"version"`
}

type ActionType int

const (
	ActionTypeTrigger ActionType = iota
	ActionTypeOperation
	ActionTypeConditional
)
