package data_test

import (
	"fmt"
	"log"
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
		os.Exit(1)
	}

	exitCode := m.Run()
	ListTables(db)

	db.Close()

	os.Exit(exitCode)
}

func ListTables(db *sqlx.DB) {
	var tables []string
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';`

	err := db.Select(&tables, query)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Tables in the 'public' schema:")
	for _, table := range tables {
		fmt.Println(table)
	}

	fmt.Println("--Next--")
}
