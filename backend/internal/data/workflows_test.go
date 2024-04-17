package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/luisya22/confluo/backend/internal/data"
	"github.com/luisya22/confluo/backend/internal/tests"
	"github.com/luisya22/confluo/backend/internal/tests/assert"
)

type workflowTestResult struct {
	workflow    data.Workflow
	shouldError bool
}

func TestWorkflowInsert(t *testing.T) {
	testMap := []struct {
		name  string
		data  data.Workflow
		wants workflowTestResult
	}{
		{
			name: "Can Insert",
			data: data.Workflow{
				UserId: tests.TestData.Users[0].Id,
				Name:   "Flow 1",
			},
			wants: workflowTestResult{
				workflow: data.Workflow{
					UserId: tests.TestData.Users[0].Id,
					Name:   "Flow 1",
				},
				shouldError: false,
			},
		},
		{
			name: "Missing UserId Should Error",
			data: data.Workflow{
				UserId: "",
				Name:   "Flow 1",
			},
			wants: workflowTestResult{
				shouldError: true,
			},
		},
		{
			name: "Missing Name Should Error",
			data: data.Workflow{
				UserId: tests.TestData.Users[0].Id,
				Name:   "",
			},
			wants: workflowTestResult{
				shouldError: true,
			},
		},
	}

	tests.SetupDb(db)
	defer tests.TeardownDb(db)

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {
			model := data.WorkflowModel{DB: db}

			err := model.Insert(&tt.data)

			if tt.wants.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NilError(t, err)

			assert.NotEqual(t, tt.data.Id, "")

			query := `SELECT * FROM workflows WHERE id = $1`

			var workflow data.Workflow

			dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err = db.QueryRowxContext(dbCtx, query, tt.data.Id).StructScan(&workflow)

			assert.NilError(t, err)

			assert.Equal(t, tt.data.UserId, workflow.UserId)
			assert.Equal(t, tt.data.Name, workflow.Name)
		})
	}
}

func TestWorkflowUpdate(t *testing.T) {
	testMap := []struct {
		name  string
		data  data.Workflow
		wants workflowTestResult
	}{
		{
			name: "Can Update",
			data: data.Workflow{
				Id:      tests.TestData.Workflows[0].Id,
				UserId:  tests.TestData.Workflows[0].UserId,
				Name:    "New Flow",
				Version: 1,
			},
			wants: workflowTestResult{
				workflow: data.Workflow{
					Id:      tests.TestData.Workflows[0].Id,
					UserId:  tests.TestData.Workflows[0].UserId,
					Name:    "New Flow",
					Version: 2,
				},
			},
		},
	}

	tests.SetupDb(db)
	defer tests.TeardownDb(db)

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {
			model := data.WorkflowModel{DB: db}

			err := model.Update(&tt.data)

			if tt.wants.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NilError(t, err)

			query := `SELECT * FROM workflows WHERE id = $1`

			var workflow data.Workflow

			dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err = db.QueryRowxContext(dbCtx, query, tt.data.Id).StructScan(&workflow)

			assert.NilError(t, err)

			assert.Equal(t, tt.data.UserId, workflow.UserId)
			assert.Equal(t, tt.data.Name, workflow.Name)
		})
	}

}
