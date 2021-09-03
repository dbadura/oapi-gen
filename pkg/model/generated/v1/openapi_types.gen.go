// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package v1

import (
	"time"
)

// Defines values for Status.
const (
	StatusDone Status = "done"

	StatusOnHold Status = "on_hold"

	StatusWorking Status = "working"
)

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Status defines model for Status.
type Status string

// Todo defines model for Todo.
type Todo struct {
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	Id          *int32     `json:"id,omitempty"`
	Status      Status     `json:"status"`
	Task        string     `json:"task"`
	User        string     `json:"user"`
}

// User defines model for User.
type User string