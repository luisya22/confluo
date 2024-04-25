package executor

import "errors"

// Interface for all the other packages
// Way to send to DB providers and actions

type Executor struct {
	providers map[string]Provider
}

type Action func(map[string]interface{}) (map[string]interface{}, error)
type Provider map[string]Action

var (
	ErrProviderNotFound = errors.New("provider not found")
	ErrActionNotFound   = errors.New("action not found")
	ErrNotTriggered     = errors.New("action not triggered")
)

func NewExecutor() *Executor {

	return &Executor{
		providers: make(map[string]Provider),
	}
}

func (e *Executor) Subscribe(name string, p Provider) error {
	e.providers[name] = p

	return nil
}

func (e *Executor) Execute(provider string, action string, params map[string]interface{}) (map[string]interface{}, error) {
	p, ok := e.providers[provider]
	if !ok {
		return params, ErrProviderNotFound
	}

	a, ok := p[action]
	if !ok {
		return params, ErrActionNotFound
	}

	return a(params)
}
