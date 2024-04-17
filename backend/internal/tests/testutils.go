package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/luisya22/confluo/backend/internal/data"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/testcontainers/testcontainers-go/wait"
)

// TestData struct that includes all other structs
var TestData = struct {
	Users           []data.User
	Actions         []data.Action
	Workflows       []data.Workflow
	WorkflowActions []data.WorkflowAction
}{
	Users: []data.User{
		{Id: "550e8400-e29b-41d4-a716-446655440000"},
		{Id: "550e8400-e29b-41d4-a716-446655440006"},
	},
	Actions: []data.Action{
		{Id: "550e8400-e29b-41d4-a716-446655440001", Operation: "Create", Provider: "System", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
		{Id: "550e8400-e29b-41d4-a716-446655440004", Operation: "Update", Provider: "System", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
		{Id: "550e8400-e29b-41d4-a716-446655440007", Operation: "Review", Provider: "System", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
		{Id: "550e8400-e29b-41d4-a716-446655440008", Operation: "Approve", Provider: "System", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
	},
	Workflows: []data.Workflow{
		{Id: "550e8400-e29b-41d4-a716-446655440002", UserId: "550e8400-e29b-41d4-a716-446655440000", Name: "User Onboarding", TriggerId: sql.NullString{String: "550e8400-e29b-41d4-a716-446655440003", Valid: true}, CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
		{Id: "550e8400-e29b-41d4-a716-446655440009", UserId: "550e8400-e29b-41d4-a716-446655440006", Name: "Document Approval", TriggerId: sql.NullString{String: "550e8400-e29b-41d4-a716-446655440010", Valid: true}, CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
	},
	WorkflowActions: []data.WorkflowAction{
		{Id: "550e8400-e29b-41d4-a716-446655440003", Text: "Begin Onboarding", Type: "Init", Params: map[string]interface{}{}, WorkflowId: "550e8400-e29b-41d4-a716-446655440002", ActionId: "550e8400-e29b-41d4-a716-446655440001", NextActionId: "550e8400-e29b-41d4-a716-446655440005", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Id: "550e8400-e29b-41d4-a716-446655440005", Text: "Complete Onboarding", Type: "Finish", Params: map[string]interface{}{}, WorkflowId: "550e8400-e29b-41d4-a716-446655440002", ActionId: "550e8400-e29b-41d4-a716-446655440004", NextActionId: "", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Id: "550e8400-e29b-41d4-a716-446655440010", Text: "Submit Document", Type: "Init", Params: map[string]interface{}{}, WorkflowId: "550e8400-e29b-41d4-a716-446655440009", ActionId: "550e8400-e29b-41d4-a716-446655440007", NextActionId: "550e8400-e29b-41d4-a716-446655440011", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Id: "550e8400-e29b-41d4-a716-446655440011", Text: "Approve Document", Type: "Finish", Params: map[string]interface{}{}, WorkflowId: "550e8400-e29b-41d4-a716-446655440009", ActionId: "550e8400-e29b-41d4-a716-446655440008", NextActionId: "", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	},
}

func NewTestDB() (*sqlx.DB, error) {
	containerReq := tc.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "pass",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		),
	}

	dbContainer, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		},
	)
	if err != nil {
		return &sqlx.DB{}, err
	}

	hostPort, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return &sqlx.DB{}, err
	}

	postgresURLTemplate := "postgres://user:pass@localhost:%s?sslmode=disable"
	postgresURL := fmt.Sprintf(postgresURLTemplate, hostPort.Port())

	db, err := sqlx.Open("postgres", postgresURL)
	if err != nil {
		return &sqlx.DB{}, err
	}

	err = SetupDb(db)
	if err != nil {
		return db, err
	}

	return db, nil

}

func SetupDb(db *sqlx.DB) error {
	script, err := os.ReadFile("../tests/testdata/setup.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(script))
	if err != nil {
		return err
	}

	fmt.Println("Hi")

	return nil
}

func TeardownDb(db *sqlx.DB) error {

	script, err := os.ReadFile("../tests/testdata/teardown.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(script))
	if err != nil {
		return err
	}

	return nil
}

func ptrToString(s string) *string {
	return &s
}

func ptrToFloat64(f float64) *float64 {
	return &f
}
