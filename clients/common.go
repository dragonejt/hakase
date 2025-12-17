// Package clients provides interfaces and types for interacting with the backend API and Discord session.
package clients

import (
	"net/http"

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
	CreateAssignment(span *sentry.Span, assignment Assignment) (string, error)
	UpdateAssignment(span *sentry.Span, assignment Assignment) error
	DeleteAssignment(span *sentry.Span, assignmentID string) error
}

type BackendClient struct {
	HakaseClient

	Url        string
	AuthToken  string
	HTTPClient *http.Client
}

type RequestOptions struct {
	Mode        string `json:"mode,omitempty"`
	ReturnEdits string `json:"returnEdits,omitempty"`
}

type ActionResponse struct {
	Validation ResponseValidation `json:"validation"`
	Edits      ResponseEdits      `json:"edits"`
}

type ResponseValidation struct {
	Result             string                     `json:"result,omitempty"`
	SubmissionCriteria []string                   `json:"submissionCriteria,omitempty"`
	Parameters         map[string]ParameterResult `json:"parameters,omitempty"`
}

type ParameterResult struct {
	Result               string   `json:"result,omitempty"`
	EvaluatedConstraints []string `json:"evaluatedConstraints,omitempty"`
	Required             bool     `json:"required,omitempty"`
}

type ResponseEdits struct {
	Type                 string         `json:"type,omitempty"`
	Edits                []ResponseEdit `json:"edits,omitempty"`
	AddedObjectCount     int            `json:"addedObjectCount,omitempty"`
	ModifiedObjectsCount int            `json:"modifiedObjectsCount,omitempty"`
	DeletedObjectsCount  int            `json:"deletedObjectsCount,omitempty"`
	AddedLinksCount      int            `json:"addedLinksCount,omitempty"`
	DeletedLinksCount    int            `json:"deletedLinksCount,omitempty"`
}

type ResponseEdit struct {
	Type       string `json:"type,omitempty"`
	PrimaryKey string `json:"primaryKey,omitempty"`
	ObjectType string `json:"objectType,omitempty"`
}

type DiscordSession struct{}
