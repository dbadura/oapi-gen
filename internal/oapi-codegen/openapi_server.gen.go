// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// Defines values for TodoStatus.
const (
	TodoStatusDone TodoStatus = "done"

	TodoStatusOnHold TodoStatus = "on_hold"

	TodoStatusWorking TodoStatus = "working"
)

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Todo defines model for Todo.
type Todo struct {
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	Id          *int32     `json:"id,omitempty"`
	Status      TodoStatus `json:"status"`
	Task        string     `json:"task"`
	User        string     `json:"user"`
}

// TodoStatus defines model for Todo.Status.
type TodoStatus string

// GetTodosParams defines parameters for GetTodos.
type GetTodosParams struct {
	// ID of the user
	User string `json:"user"`

	// filter todo's by status
	Status *GetTodosParamsStatus `json:"status,omitempty"`
}

// GetTodosParamsStatus defines parameters for GetTodos.
type GetTodosParamsStatus string

// CreateTodoJSONBody defines parameters for CreateTodo.
type CreateTodoJSONBody Todo

// CreateTodoJSONRequestBody defines body for CreateTodo for application/json ContentType.
type CreateTodoJSONRequestBody CreateTodoJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /todos)
	GetTodos(w http.ResponseWriter, r *http.Request, params GetTodosParams)

	// (POST /todos)
	CreateTodo(w http.ResponseWriter, r *http.Request)

	// (DELETE /todos/{todoId})
	DeleteTodo(w http.ResponseWriter, r *http.Request, todoId int32)

	// (PUT /todos/{todoId})
	UpdateTodo(w http.ResponseWriter, r *http.Request, todoId int32)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// GetTodos operation middleware
func (siw *ServerInterfaceWrapper) GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTodosParams

	// ------------- Required query parameter "user" -------------
	if paramValue := r.URL.Query().Get("user"); paramValue != "" {

	} else {
		http.Error(w, "Query argument user is required, but not found", http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "user", r.URL.Query(), &params.User)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter user: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "status" -------------
	if paramValue := r.URL.Query().Get("status"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "status", r.URL.Query(), &params.Status)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter status: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTodos(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// CreateTodo operation middleware
func (siw *ServerInterfaceWrapper) CreateTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateTodo(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// DeleteTodo operation middleware
func (siw *ServerInterfaceWrapper) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "todoId" -------------
	var todoId int32

	err = runtime.BindStyledParameter("simple", false, "todoId", chi.URLParam(r, "todoId"), &todoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter todoId: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteTodo(w, r, todoId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// UpdateTodo operation middleware
func (siw *ServerInterfaceWrapper) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "todoId" -------------
	var todoId int32

	err = runtime.BindStyledParameter("simple", false, "todoId", chi.URLParam(r, "todoId"), &todoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter todoId: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateTodo(w, r, todoId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL     string
	BaseRouter  chi.Router
	Middlewares []MiddlewareFunc
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/todos", wrapper.GetTodos)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/todos", wrapper.CreateTodo)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/todos/{todoId}", wrapper.DeleteTodo)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/todos/{todoId}", wrapper.UpdateTodo)
	})

	return r
}
