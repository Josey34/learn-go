package domain

import "fmt"

type ValidationError struct {
	Field   string // Which field failed (e.g., "title")
	Message string // What went wrong (e.g., "is required")
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s", e.Field, e.Message)
}

type NotFoundError struct {
	Resource string // The resource that was not found (e.g., "Task")
	ID       int    // The ID of the resource
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

type DatabaseError struct {
	Operation string // The database operation that failed (e.g., "insert", "update")
	Err       error  // The underlying error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Database error during %s: %v", e.Operation, e.Err)
}

type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type QueueError struct {
	Message string
}

func (e *QueueError) Error() string {
	return e.Message
}

var (
	ErrQueueFull  = &QueueError{Message: "queue is full"}
	ErrQueueEmpty = &QueueError{Message: "queue is empty"}
)
