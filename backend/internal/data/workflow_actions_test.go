package data_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/luisya22/confluo/backend/internal/data"
	"github.com/luisya22/confluo/backend/internal/tests"
	"github.com/luisya22/confluo/backend/internal/tests/assert"
)

type workflowActionTestResult struct {
	workflowAction data.WorkflowAction
	shouldError    bool
}

func TestInsertWorkflowAction(t *testing.T) {
	testMap := []struct {
		name  string
		data  data.WorkflowAction
		wants workflowActionTestResult
	}{
		{
			name: "Can Insert",
			data: data.WorkflowAction{
				Text:       "Action 1",
				WorkflowId: tests.Data.Workflows[0].Id,
				ActionId:   tests.Data.Actions[0].Id,
				Type:       "operation",
			},
			wants: workflowActionTestResult{
				workflowAction: data.WorkflowAction{
					Text:       "Action 1",
					WorkflowId: tests.Data.Workflows[0].Id,
					ActionId:   tests.Data.Actions[0].Id,
					Type:       "operation",
				},

				shouldError: false,
			},
		},
		{
			name: "WorkflowId Missing Should Error",
			data: data.WorkflowAction{
				Text:     "Action 1",
				ActionId: tests.Data.Actions[0].Id,
				Type:     "operation",
			},
			wants: workflowActionTestResult{
				shouldError: true,
			},
		},
		{
			name: "WorkflowId Missing Should Error",
			data: data.WorkflowAction{
				Text:     "Action 1",
				ActionId: tests.Data.Actions[0].Id,
				Type:     "operation",
			},
			wants: workflowActionTestResult{
				shouldError: true,
			},
		},
	}

	tests.SetupDb(db)
	defer tests.TeardownDb(db)

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {
			model := data.WorkFlowActionModel{DB: db}

			err := model.Insert(&tt.data)

			if tt.wants.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NilError(t, err)

			query := `SELECT id, text, workflow_id, action_id, type, next_action_id
				FROM workflow_actions WHERE id = $1`

			var workflowAction data.WorkflowAction

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err = db.QueryRowxContext(ctx, query, tt.data.Id).Scan(
				&workflowAction.Id,
				&workflowAction.Text,
				&workflowAction.WorkflowId,
				&workflowAction.ActionId,
				&workflowAction.Type,
				&workflowAction.NextActionId,
			)

			assert.NilError(t, err)

			assert.NotEqual(t, tt.data.Id, "")
			assert.Equal(t, tt.data.Text, tt.wants.workflowAction.Text)
			assert.Equal(t, tt.data.WorkflowId, tt.wants.workflowAction.WorkflowId)
			assert.Equal(t, tt.data.ActionId, tt.wants.workflowAction.ActionId)
			assert.Equal(t, tt.data.Type, tt.wants.workflowAction.Type)
			assert.Equal(t, tt.data.NextActionId, tt.wants.workflowAction.NextActionId)
		})
	}
}

func TestWorkflowActionUpdate(t *testing.T) {
	testMap := []struct {
		name  string
		data  data.WorkflowAction
		wants workflowActionTestResult
	}{
		{
			name: "Can Update",
			data: data.WorkflowAction{
				Id:   tests.Data.WorkflowActions[0].Id,
				Text: "New Text",
				Params: map[string]interface{}{
					"param1": "string1",
					"param2": 2,
				},
				Version: 1,
			},
			wants: workflowActionTestResult{
				workflowAction: data.WorkflowAction{
					Id:   tests.Data.WorkflowActions[0].Id,
					Text: "New Text",
					Params: map[string]interface{}{
						"param1": "string1",
						"param2": 2,
					},
					Version: 1,
				},
				shouldError: false,
			},
		},
	}

	err := tests.SetupDb(db)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer tests.TeardownDb(db)

	ListTables(db)

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {
			model := data.WorkFlowActionModel{DB: db}

			err := model.Update(&tt.data)

			if tt.wants.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NilError(t, err)

			query := `SELECT id, text, workflow_id, action_id, type, next_action_id, params
				FROM workflow_actions WHERE id = $1`

			var workflowAction data.WorkflowAction

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var rawData []byte
			err = db.QueryRowxContext(ctx, query, tt.data.Id).Scan(
				&workflowAction.Id,
				&workflowAction.Text,
				&workflowAction.WorkflowId,
				&workflowAction.ActionId,
				&workflowAction.Type,
				&workflowAction.NextActionId,
				&rawData,
			)

			var params map[string]interface{}
			if rawData != nil {
				err := json.Unmarshal(rawData, &params)
				if err != nil {
					t.Fatal(err)
				}
			}

			assert.NilError(t, err)

			assert.NotEqual(t, tt.data.Id, "")
			assert.Equal(t, tt.data.Text, tt.wants.workflowAction.Text)
			assert.Equal(t, tt.data.WorkflowId, tt.wants.workflowAction.WorkflowId)
			assert.Equal(t, tt.data.ActionId, tt.wants.workflowAction.ActionId)
			assert.Equal(t, tt.data.Type, tt.wants.workflowAction.Type)
			assert.Equal(t, tt.data.NextActionId, tt.wants.workflowAction.NextActionId)

			for k, v := range workflowAction.Params {
				wants := tt.wants.workflowAction.Params[k]
				assert.Equal(t, v, wants)
			}

		})
	}
}
