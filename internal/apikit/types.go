package api

import "net/http"

var contentTypesForFiles = []string{"application/json", "image/png", "image/jpeg", "image/tiff", "image/webp", "image/svg+xml", "image/gif", "image/tiff", "image/x-icon", "application/pdf", "application/octet-stream"}

type GetTodosRequest struct {
	User   string `validate:"regex1"`
	Status *string
}

type GetTodosResponse interface {
	isGetTodosResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// get todo's response
type GetTodos200Response struct{}

func (r *GetTodos200Response) isGetTodosResponse() {}

func (r *GetTodos200Response) StatusCode() int {
	return 200
}

func (r *GetTodos200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

type CreateTodoRequest struct{}

type CreateTodoResponse interface {
	isCreateTodoResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Todo creation response
type CreateTodo201Response struct{}

func (r *CreateTodo201Response) isCreateTodoResponse() {}

func (r *CreateTodo201Response) StatusCode() int {
	return 201
}

func (r *CreateTodo201Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(201)
	return nil
}

type DeleteTodoRequest struct {
	TodoId int32
}

type DeleteTodoResponse interface {
	isDeleteTodoResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// no content
type DeleteTodo204Response struct{}

func (r *DeleteTodo204Response) isDeleteTodoResponse() {}

func (r *DeleteTodo204Response) StatusCode() int {
	return 204
}

func (r *DeleteTodo204Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(204)
	return nil
}

type UpdateTodoRequest struct {
	TodoId int32
}

type UpdateTodoResponse interface {
	isUpdateTodoResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

type UpdateTodo200Response struct{}

func (r *UpdateTodo200Response) isUpdateTodoResponse() {}

func (r *UpdateTodo200Response) StatusCode() int {
	return 200
}

func (r *UpdateTodo200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}
