package data

import "github.com/jmoiron/sqlx"

type User struct {
	Id string
}

type UserModle struct {
	DB *sqlx.DB
}
