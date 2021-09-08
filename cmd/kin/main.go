package main

import (
	"context"
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/legacy"
	"log"
	"net/http"
	"os"
)

func main() {
	openapi3.DefineStringFormat("uuid", openapi3.FormatOfStringForUUIDOfRFC4122)

	ctx := context.Background()
	loader := &openapi3.Loader{}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Printf("Current dir: %s", pwd)
	doc, err := loader.LoadFromFile("./openapi/basic.yaml")
	if err != nil {
		panic(err)
	}

	err = doc.Validate(ctx)
	if err != nil {
		panic(err)
	}
	router, _ := legacy.NewRouter(doc)
	httpReq, _ := http.NewRequest(http.MethodGet, "http://localhost:32208/todos?user=0d5db390-0f8d-4e81-8bcc-2a326483c761", nil)

	// Find route
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		panic(err)
	}
	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		panic(err)
	}

	resp := map[string]interface{}{
		"id": 0,
		"user": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"task": "string",
		"status": "on_hold",
		"created_at": "2021-09-08T12:36:15.437Z",
		"completed_at": "2021-09-08T12:36:15.437Z",
	}

	//out, err :=json.Marshal([]map[string]interface{}{resp})
	//if err != nil {
	//	panic(err)
	//}
	var (
		respStatus      = 200
		respContentType = "application/json"
		respBody        = []map[string]interface{}{resp}
	)

	log.Println("Response")
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header:                 http.Header{"Content-Type": []string{respContentType}},
	}
	data, _ := json.Marshal(respBody)
	responseValidationInput.SetBodyBytes(data)

	// Validate response.
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}
}
