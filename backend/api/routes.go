package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	// routes
	// Get Workflows
	// Get Workflow
	// Post Workflow
	// Patch Workflow
	// Post WorkflowAction
	// Patch WorkflowAction

	// Github Oauth
	// Google Sheets Oauth

	// Authentication
	router.Group(func(r chi.Router) {
		r.Post("/auth/github/callback")
	})

	return router

}
