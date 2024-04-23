package data_test

import (
	"context"
	"math"
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
				UserId: tests.Data.Users[0].Id,
				Name:   "Flow 1",
			},
			wants: workflowTestResult{
				workflow: data.Workflow{
					UserId: tests.Data.Users[0].Id,
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
				UserId: tests.Data.Users[0].Id,
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
				Id:      tests.Data.Workflows[0].Id,
				UserId:  tests.Data.Workflows[0].UserId,
				Name:    "New Flow",
				Version: 1,
			},
			wants: workflowTestResult{
				workflow: data.Workflow{
					Id:      tests.Data.Workflows[0].Id,
					UserId:  tests.Data.Workflows[0].UserId,
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

func TestWorkflowGet(t *testing.T) {

	testMap := []struct {
		name  string
		data  string
		wants workflowTestResult
	}{
		{
			name: "Can Get",
			data: tests.Data.Workflows[0].Id,
			wants: workflowTestResult{
				workflow:    tests.Data.Workflows[0],
				shouldError: false,
			},
		},
		{
			name: "Wrong Id Should Error",
			data: "00000000-0000-0000-0000-000000000000",
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

			workflow, err := model.Get(tt.data)

			if tt.wants.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NilError(t, err)

			assert.Equal(t, workflow.Id, tt.wants.workflow.Id)
			assert.Equal(t, workflow.UserId, tt.wants.workflow.UserId)
			assert.Equal(t, workflow.Name, tt.wants.workflow.Name)
			assert.Equal(t, workflow.TriggerId.String, tt.wants.workflow.TriggerId.String)

			for i, a := range workflow.Actions {
				assert.Equal(t, a.Id, tt.wants.workflow.Actions[i].Id)
				assert.Equal(t, a.Text, tt.wants.workflow.Actions[i].Text)
				assert.Equal(t, a.WorkflowId, tt.wants.workflow.Actions[i].WorkflowId)
				assert.Equal(t, a.ActionId, tt.wants.workflow.Actions[i].ActionId)
				assert.Equal(t, a.Type, tt.wants.workflow.Actions[i].Type)
				assert.Equal(t, a.NextActionId, tt.wants.workflow.Actions[i].NextActionId)

				assert.Equal(t, a.Action.Id, tt.wants.workflow.Actions[i].Action.Id)
				assert.Equal(t, a.Action.Operation, tt.wants.workflow.Actions[i].Action.Operation)
				assert.Equal(t, a.Action.Provider.Id, tt.wants.workflow.Actions[i].Action.Provider.Id)
				assert.Equal(t, a.Action.Provider.Name, tt.wants.workflow.Actions[i].Action.Provider.Name)

			}
		})
	}
}

func TestWorkflowDelete(t *testing.T) {
	testMap := []struct {
		name  string
		data  string
		wants workflowTestResult
	}{
		{
			name: "Can Delete",
			data: tests.Data.Workflows[0].Id,
			wants: workflowTestResult{

				shouldError: false,
			},
		},
	}

	tests.SetupDb(db)
	defer tests.TeardownDb(db)

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {
			model := data.WorkflowModel{DB: db}

			err := model.Delete(tt.data)

			if tt.wants.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NilError(t, err)

			query := `SELECT * FROM workflows WHERE id = $1`

			var workflow data.Workflow

			dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err = db.QueryRowxContext(dbCtx, query, tt.data).StructScan(&workflow)

			assert.Error(t, err)
		})
	}
}

func TestWorkflowGetAll(t *testing.T) {
	type params struct {
		userId  string
		filters data.Filters
	}

	type getAllWorkflowsTestResult struct {
		metadata    data.Metadata
		shouldError bool
	}

	testMap := []struct {
		name  string
		data  params
		wants getAllWorkflowsTestResult
	}{}

	tests.SetupDb(db)
	defer tests.TeardownDb(db)

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {
			model := data.WorkflowModel{DB: db}

			workflows, metadata, err := model.GetAll(tt.data.userId, tt.data.filters)

			if tt.wants.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NilError(t, err)

			maxIndex := int(math.Min(float64(tt.wants.metadata.PageSize), float64(len(tests.Data.Workflows))))
			var wantsWorkflows []data.Workflow = tests.Data.Workflows

			wantsWorkflows = wantsWorkflows[:maxIndex]

			assert.Equal(t, len(workflows), len(wantsWorkflows))
			for workflowIndex, workflow := range workflows {
				assert.Equal(t, workflow.Id, wantsWorkflows[workflowIndex].Id)
				assert.Equal(t, workflow.UserId, wantsWorkflows[workflowIndex].UserId)
				assert.Equal(t, workflow.Name, wantsWorkflows[workflowIndex].Name)
				assert.Equal(t, workflow.TriggerId.String, wantsWorkflows[workflowIndex].TriggerId.String)

				for actionIndex, a := range workflow.Actions {
					assert.Equal(t, a.Id, wantsWorkflows[workflowIndex].Actions[actionIndex].Id)
					assert.Equal(t, a.Text, wantsWorkflows[workflowIndex].Actions[actionIndex].Text)
					assert.Equal(t, a.WorkflowId, wantsWorkflows[workflowIndex].Actions[actionIndex].WorkflowId)
					assert.Equal(t, a.ActionId, wantsWorkflows[workflowIndex].Actions[actionIndex].ActionId)
					assert.Equal(t, a.Type, wantsWorkflows[workflowIndex].Actions[actionIndex].Type)
					assert.Equal(t, a.NextActionId, wantsWorkflows[workflowIndex].Actions[actionIndex].NextActionId)

					assert.Equal(t, a.Action.Id, wantsWorkflows[workflowIndex].Actions[actionIndex].Action.Id)
					assert.Equal(t, a.Action.Operation, wantsWorkflows[workflowIndex].Actions[actionIndex].Action.Operation)
					assert.Equal(t, a.Action.Provider.Id, wantsWorkflows[workflowIndex].Actions[actionIndex].Action.Provider.Id)
					assert.Equal(t, a.Action.Provider.Name, wantsWorkflows[workflowIndex].Actions[actionIndex].Action.Provider.Name)

				}
			}

			assert.Equal(t, metadata.CurrentPage, tt.wants.metadata.CurrentPage)
			assert.Equal(t, metadata.PageSize, tt.wants.metadata.PageSize)
			assert.Equal(t, metadata.FirstPage, tt.wants.metadata.FirstPage)
			assert.Equal(t, metadata.LastPage, tt.wants.metadata.LastPage)
			assert.Equal(t, metadata.TotalRecords, tt.wants.metadata.TotalRecords)

		})
	}
}
