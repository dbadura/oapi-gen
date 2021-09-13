package api

import (
	"context"
	"fmt"
	routing "github.com/go-ozzo/ozzo-routing"
	validator "github.com/go-playground/validator"
	"net/http"
	"regexp"
)

func NewTodoAppOASServer(options *ServerOpts) *TodoAppOASServer {
	serverWrapper := &TodoAppOASServer{Server: newServer(options), Validator: NewValidation()}
	serverWrapper.Server.SwaggerSpec = swagger
	serverWrapper.registerValidators()
	return serverWrapper
}

type TodoAppOASServer struct {
	*Server
	Validator         *Validator
	getTodosHandler   *getTodosHandlerRoute
	createTodoHandler *createTodoHandlerRoute
	deleteTodoHandler *deleteTodoHandlerRoute
	updateTodoHandler *updateTodoHandlerRoute
}

// Returns all the todo's of the user
type GetTodosHandler func(ctx context.Context, request *GetTodosRequest) GetTodosResponse

type getTodosHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetTodosHandler
}

func (server *TodoAppOASServer) SetGetTodosHandler(handler GetTodosHandler, middleware ...Middleware) {
	server.getTodosHandler = &getTodosHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/todos", Handler: server.GetTodosHandler, Middleware: middleware}}
}

func (server *TodoAppOASServer) GetTodosHandler(c *routing.Context) error {
	if server.getTodosHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetTodos (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetTodosRequest)
		if len(c.Request.URL.Query()["user"]) > 0 {
			if err := fromString(c.Request.URL.Query()["user"][0], &request.User); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetTodos (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.URL.Query()["status"]) > 0 {
			if err := fromString(c.Request.URL.Query()["status"][0], &request.Status); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetTodos (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetTodos (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getTodosHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetTodos (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetTodos (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

// Creates a new todo
type CreateTodoHandler func(ctx context.Context, request *CreateTodoRequest) CreateTodoResponse

type createTodoHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    CreateTodoHandler
}

func (server *TodoAppOASServer) SetCreateTodoHandler(handler CreateTodoHandler, middleware ...Middleware) {
	server.createTodoHandler = &createTodoHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "POST", Path: "/todos", Handler: server.CreateTodoHandler, Middleware: middleware}}
}

func (server *TodoAppOASServer) CreateTodoHandler(c *routing.Context) error {
	if server.createTodoHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: CreateTodo (POST) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(CreateTodoRequest)
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateTodo (POST) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.createTodoHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: CreateTodo (POST) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateTodo (POST) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

// delete a todo
type DeleteTodoHandler func(ctx context.Context, request *DeleteTodoRequest) DeleteTodoResponse

type deleteTodoHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DeleteTodoHandler
}

func (server *TodoAppOASServer) SetDeleteTodoHandler(handler DeleteTodoHandler, middleware ...Middleware) {
	server.deleteTodoHandler = &deleteTodoHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "DELETE", Path: "/todos/<todoId>", Handler: server.DeleteTodoHandler, Middleware: middleware}}
}

func (server *TodoAppOASServer) DeleteTodoHandler(c *routing.Context) error {
	if server.deleteTodoHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DeleteTodo (DELETE) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DeleteTodoRequest)
		if err := fromString(c.Param("todoId"), &request.TodoId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteTodo (DELETE) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteTodo (DELETE) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.deleteTodoHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DeleteTodo (DELETE) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteTodo (DELETE) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

// Updates the status of a todo
type UpdateTodoHandler func(ctx context.Context, request *UpdateTodoRequest) UpdateTodoResponse

type updateTodoHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    UpdateTodoHandler
}

func (server *TodoAppOASServer) SetUpdateTodoHandler(handler UpdateTodoHandler, middleware ...Middleware) {
	server.updateTodoHandler = &updateTodoHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "PUT", Path: "/todos/<todoId>", Handler: server.UpdateTodoHandler, Middleware: middleware}}
}

func (server *TodoAppOASServer) UpdateTodoHandler(c *routing.Context) error {
	if server.updateTodoHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: UpdateTodo (PUT) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(UpdateTodoRequest)
		if err := fromString(c.Param("todoId"), &request.TodoId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: UpdateTodo (PUT) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: UpdateTodo (PUT) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.updateTodoHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: UpdateTodo (PUT) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: UpdateTodo (PUT) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

func (server *TodoAppOASServer) registerValidators() {
	regex1 := regexp.MustCompile("^([a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12})?$")
	callbackRegex1 := func(fl validator.FieldLevel) bool {
		return regex1.MatchString(fl.Field().String())
	}
	server.Validator.RegisterValidation("regex1", callbackRegex1)
}

func (server *TodoAppOASServer) Start(port int) error {
	routes := []RouteDescription{}
	if server.getTodosHandler != nil {
		routes = append(routes, server.getTodosHandler.routeDescription)
	}
	if server.createTodoHandler != nil {
		routes = append(routes, server.createTodoHandler.routeDescription)
	}
	if server.deleteTodoHandler != nil {
		routes = append(routes, server.deleteTodoHandler.routeDescription)
	}
	if server.updateTodoHandler != nil {
		routes = append(routes, server.updateTodoHandler.routeDescription)
	}
	return server.Server.Start(port, routes)
}

const swagger = "{\"info\":{\"description\":\"OpenApi specification for a todo application\",\"title\":\"Todo app OAS\",\"version\":\"1.0.0\"},\"paths\":{\"/todos\":{\"get\":{\"description\":\"Returns all the todo's of the user\",\"operationId\":\"GetTodos\",\"parameters\":[{\"description\":\"ID of the user\",\"name\":\"user\",\"in\":\"query\",\"required\":true,\"schema\":{\"type\":\"string\",\"format\":\"uuid\"}},{\"description\":\"filter todo's by status\",\"name\":\"status\",\"in\":\"query\",\"schema\":{\"type\":\"string\",\"enum\":[\"on_hold\",\"working\",\"done\"]}}],\"responses\":{\"200\":{\"description\":\"get todo's response\"},\"default\":{\"description\":\"unexpected Error\"}}},\"post\":{\"description\":\"Creates a new todo\",\"operationId\":\"CreateTodo\",\"responses\":{\"201\":{\"description\":\"Todo creation response\"},\"default\":{\"description\":\"unexpected error\"}}}},\"/todos/{todoId}\":{\"put\":{\"description\":\"Updates the status of a todo\",\"operationId\":\"UpdateTodo\",\"parameters\":[{\"description\":\"Id of the todo\",\"name\":\"todoId\",\"in\":\"path\",\"required\":true,\"schema\":{\"type\":\"integer\",\"format\":\"int32\"}}],\"responses\":{\"200\":{},\"default\":{\"description\":\"unexpected error\"}}},\"delete\":{\"description\":\"delete a todo\",\"operationId\":\"DeleteTodo\",\"parameters\":[{\"description\":\"Id of the todo\",\"name\":\"todoId\",\"in\":\"path\",\"required\":true,\"schema\":{\"type\":\"integer\",\"format\":\"int32\"}}],\"responses\":{\"204\":{\"description\":\"no content\"},\"default\":{\"description\":\"unexpected error\"}}}}}}"
