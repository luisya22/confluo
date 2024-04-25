package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v61/github"
	"github.com/luisya22/confluo/backend/internal/executor"
)

const ProviderName = "Github"

func Initialize(e *executor.Executor) {
	actions := make(executor.Provider)

	actions["New Issue"] = newIssue
	actions["New Branch"] = newBranch
	actions["New Commit"] = newCommit
	actions["New Repo"] = newRepo
	actions["New Release"] = newRelease

	actions["Create Comment"] = createComment
	actions["Create Issue"] = createIssue
	actions["Update Issue"] = updateIssue
	actions["Create Pull Request"] = createPullRequest
	actions["Update Pull Request"] = updatePullRequest
	actions["Delete Branch"] = deleteBranch

	actions["Find Issue"] = findIssue
	actions["Find Pull Request"] = findPullRequest

	e.Subscribe(ProviderName, actions)
}

// Triggers

func newIssue(params map[string]interface{}) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, owner, repo, err := getRepoData(params)
	if err != nil {
		return params, err
	}

	lastIssue, ok := params["lastIssue"].(int)
	if !ok {
		return params, fmt.Errorf("lastIssue not found or it is not correct format")
	}

	client := github.NewClient(nil).WithAuthToken(token)

	for {
		lastIssue++

		issue, res, err := client.Issues.Get(ctx, owner, repo, lastIssue)
		if err != nil {
			if res != nil && res.StatusCode == http.StatusNotFound {
				return params, executor.ErrNotTriggered
			}

			return params, err
		}

		if !issue.IsPullRequest() {
			// TODO: Finish this
			params["issueTitle"] = *issue.Title
			params["lastIssue"] = *issue.Number

			fmt.Println(*issue.Title)
			break
		}

	}

	return params, nil
}

func newBranch(params map[string]interface{}) (map[string]interface{}, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	//
	// token, owner, repo, err := getRepoData(params)
	// if err != nil {
	// 	return params, err
	// }

	return params, fmt.Errorf("not implemented")
}

func newCommit(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

func newRepo(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

func newRelease(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

// Events

func createComment(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

func createIssue(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

func updateIssue(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

func createPullRequest(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

func updatePullRequest(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

func deleteBranch(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

func findIssue(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

func findPullRequest(params map[string]interface{}) (map[string]interface{}, error) {
	return params, fmt.Errorf("not implemented")
}

// Get Params and returns token, owner, repo and if its error
func getRepoData(params map[string]interface{}) (string, string, string, error) {
	token, ok := params["token"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("token not found or it is not correct format")
	}

	owner, ok := params["owner"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("owner not found or it is not correct format")
	}

	repo, ok := params["repo"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("repo not found or it is not correct format")
	}

	return token, owner, repo, nil

}
