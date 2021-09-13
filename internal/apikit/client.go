package api

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type TodoAppOASClient interface {
	GetTodos(request *GetTodosRequest) (GetTodosResponse, error)
	CreateTodo(request *CreateTodoRequest) (CreateTodoResponse, error)
	DeleteTodo(request *DeleteTodoRequest) (DeleteTodoResponse, error)
	UpdateTodo(request *UpdateTodoRequest) (UpdateTodoResponse, error)
}

func NewTodoAppOASClient(httpClient *http.Client, baseUrl string, options Opts) TodoAppOASClient {
	return &todoAppOASClient{httpClient: newHttpClientWrapper(httpClient, baseUrl), baseURL: baseUrl, hooks: options.Hooks, ctx: options.Ctx, xmlMatcher: regexp.MustCompile("^application\\/(.+)xml$")}
}

type todoAppOASClient struct {
	baseURL    string
	hooks      HooksClient
	ctx        context.Context
	httpClient *httpClientWrapper
	xmlMatcher *regexp.Regexp
}

// Returns all the todo's of the user
func (client *todoAppOASClient) GetTodos(request *GetTodosRequest) (GetTodosResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/todos"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	query := make(url.Values)
	query.Add("user", toString(request.User))
	if request.Status != nil {
		query.Add("status", toString(request.Status))
	}
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	// set all headers from client context
	err := setRequestHeadersFromContext(httpContext, httpRequest.Header)
	if err != nil {
		return nil, err
	}
	if len(httpRequest.Header["accept"]) == 0 && len(httpRequest.Header["Accept"]) == 0 {
		httpRequest.Header["Accept"] = []string{"application/json"}
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetTodos200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

// Creates a new todo
func (client *todoAppOASClient) CreateTodo(request *CreateTodoRequest) (CreateTodoResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/todos"
	method := "POST"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	// set all headers from client context
	err := setRequestHeadersFromContext(httpContext, httpRequest.Header)
	if err != nil {
		return nil, err
	}
	if len(httpRequest.Header["accept"]) == 0 && len(httpRequest.Header["Accept"]) == 0 {
		httpRequest.Header["Accept"] = []string{"application/json"}
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusCreated {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateTodo201Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

// delete a todo
func (client *todoAppOASClient) DeleteTodo(request *DeleteTodoRequest) (DeleteTodoResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/todos/{todoId}"
	method := "DELETE"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{todoId}", url.QueryEscape(toString(request.TodoId)), 1)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	// set all headers from client context
	err := setRequestHeadersFromContext(httpContext, httpRequest.Header)
	if err != nil {
		return nil, err
	}
	if len(httpRequest.Header["accept"]) == 0 && len(httpRequest.Header["Accept"]) == 0 {
		httpRequest.Header["Accept"] = []string{"application/json"}
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusNoContent {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteTodo204Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

// Updates the status of a todo
func (client *todoAppOASClient) UpdateTodo(request *UpdateTodoRequest) (UpdateTodoResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/todos/{todoId}"
	method := "PUT"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{todoId}", url.QueryEscape(toString(request.TodoId)), 1)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	// set all headers from client context
	err := setRequestHeadersFromContext(httpContext, httpRequest.Header)
	if err != nil {
		return nil, err
	}
	if len(httpRequest.Header["accept"]) == 0 && len(httpRequest.Header["Accept"]) == 0 {
		httpRequest.Header["Accept"] = []string{"application/json"}
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(UpdateTodo200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}
