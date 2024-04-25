package api

import (
	"net/http"
)

func (app *Application) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")

	accessToken, err := app.oauhtService.GetGithubAccessToken(code)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	userData, err := app.oauhtService.GetGithubData(accessToken)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
