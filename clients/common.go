// Package clients provides interfaces and types for interacting with the backend API and Discord session.
package clients

import (
	"github.com/getsentry/sentry-go"
)

type HakaseClient interface {
	// Course APIs
	ReadCourse(span *sentry.Span, courseID string) (Course, error)
	CreateCourse(span *sentry.Span, course Course) error
	UpdateCourse(span *sentry.Span, course Course) error
	DeleteCourse(span *sentry.Span, courseID string) error
	// Assignment APIs
	ReadAssignment(span *sentry.Span, assignmentID string) (Assignment, error)
	ListAssignments(span *sentry.Span, courseID string) ([]Assignment, error)
	CreateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error)
	UpdateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error)
	DeleteAssignment(span *sentry.Span, assignmentID string) error
}

type OntologyClient = ClientWithResponses

type OntologyEdit struct {
	Type       string `json:"type"`
	PrimaryKey string `json:"primaryKey"`
	ObjectType string `json:"objectType"`
}

type DiscordSession struct{}
