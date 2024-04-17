package data_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/luisya22/confluo/backend/internal/tests"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error

	db, err = tests.NewTestDB()
	if err != nil {
		fmt.Println("error creating db:", err)
	}

	exitCode := m.Run()

	db.Close()

	os.Exit(exitCode)
}
