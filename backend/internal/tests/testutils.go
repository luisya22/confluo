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

type TestData struct {
	Users           []data.User
	Actions         []data.Action
	Workflows       []data.Workflow
	WorkflowActions []data.WorkflowAction
}

// Data struct that includes all other structs
var Data TestData

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

	LoadTestData()

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

	return nil
}

func TeardownDb(db *sqlx.DB) error {
	script, err := os.ReadFile("../tests/testdata/teardown.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(script))
	if err != nil {
		fmt.Println(err)
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

func LoadTestData() {
	users := []data.User{
		{Id: "550e8400-e29b-41d4-a716-446655440000"},
		{Id: "550e8400-e29b-41d4-a716-446655440006"},
	}

	providers := []data.Provider{
		{
			Id:        "c4f9b885-2df5-4b1b-9fa4-81f87f824da8",
			Name:      "System",
			Logo:      "image.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
	}

	actions := []data.Action{
		{
			Id:         "550e8400-e29b-41d4-a716-446655440001",
			Operation:  "Create",
			ProviderId: providers[0].Id,
			Provider:   providers[0],
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Version:    1,
		},
		{
			Id:         "550e8400-e29b-41d4-a716-446655440004",
			Operation:  "Update",
			ProviderId: providers[0].Id,
			Provider:   providers[0],
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Version:    1,
		},
		{
			Id:         "550e8400-e29b-41d4-a716-446655440007",
			Operation:  "Review",
			ProviderId: providers[0].Id,
			Provider:   providers[0],
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Version:    1,
		},
		{
			Id:         "550e8400-e29b-41d4-a716-446655440008",
			Operation:  "Approve",
			ProviderId: providers[0].Id,
			Provider:   providers[0],
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Version:    1,
		},
	}

	workflows := []data.Workflow{
		{
			Id:        "550e8400-e29b-41d4-a716-446655440002",
			UserId:    users[0].Id,
			Name:      "User Onboarding",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
		{
			Id:        "550e8400-e29b-41d4-a716-446655440009",
			UserId:    "550e8400-e29b-41d4-a716-446655440006",
			Name:      "Document Approval",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
	}

	workflowActions := []data.WorkflowAction{
		{
			Id:           "550e8400-e29b-41d4-a716-446655440003",
			Text:         "Begin Onboarding",
			Type:         "Init",
			Params:       map[string]interface{}{},
			WorkflowId:   workflows[0].Id,
			ActionId:     actions[0].Id,
			Action:       actions[0],
			NextActionId: sql.NullString{String: "550e8400-e29b-41d4-a716-446655440005", Valid: true},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Id:           "550e8400-e29b-41d4-a716-446655440005",
			Text:         "Complete Onboarding",
			Type:         "Finish",
			Params:       map[string]interface{}{},
			WorkflowId:   workflows[0].Id,
			ActionId:     actions[1].Id,
			Action:       actions[1],
			NextActionId: sql.NullString{String: "", Valid: false},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Id:           "550e8400-e29b-41d4-a716-446655440010",
			Text:         "Submit Document",
			Type:         "Init",
			Params:       map[string]interface{}{},
			WorkflowId:   workflows[1].Id,
			ActionId:     actions[2].Id,
			Action:       actions[2],
			NextActionId: sql.NullString{String: "550e8400-e29b-41d4-a716-446655440011", Valid: true},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Id:           "550e8400-e29b-41d4-a716-446655440011",
			Text:         "Approve Document",
			Type:         "Finish",
			Params:       map[string]interface{}{},
			WorkflowId:   workflows[1].Id,
			ActionId:     actions[3].Id,
			Action:       actions[3],
			NextActionId: sql.NullString{String: "", Valid: false},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	workflows[0].Actions = append(workflows[0].Actions, workflowActions[0], workflowActions[1])
	workflows[1].Actions = append(workflows[1].Actions, workflowActions[2], workflowActions[3])

	workflows[0].TriggerId = sql.NullString{String: workflows[0].Actions[0].Id, Valid: true}
	workflows[1].TriggerId = sql.NullString{String: workflows[1].Actions[0].Id, Valid: true}

	Data = TestData{
		Users:           users,
		Actions:         actions,
		Workflows:       workflows,
		WorkflowActions: workflowActions,
	}
}
