// Package clients provides interfaces and types for interacting with the backend API and Discord session.
package clients

import (
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
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

type DatabaseClient struct {
	HakaseClient
	DB *gorm.DB
}

type DiscordSession struct{}

func MigrateDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&Course{})
	if err != nil {
		return stacktrace.Propagate(err, "error migrating course model")
	}
	err = db.AutoMigrate(&Assignment{})
	if err != nil {
		return stacktrace.Propagate(err, "error migrating assignment model")
	}

	return nil
}
