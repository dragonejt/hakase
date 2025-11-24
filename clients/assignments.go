// Package clients implements backend API operations for assignments.
package clients

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

type Assignment struct {
	gorm.Model
	Course Course `gorm:"foreignKey:ID;references:ID"`
	Name   string
	Due    time.Time `gorm:"column:delete_time"`
	Link   string
}

// ReadAssignment retrieves an assignment by its ID from the backend.
func (client *DatabaseClient) ReadAssignment(span *sentry.Span, assignmentID string) (Assignment, error) {
	span = span.StartChild("readAssignment")
	defer span.Finish()

	return gorm.G[Assignment](client.DB).Where("id = ?", assignmentID).First(span.Context())
}

// ListAssignments lists all assignments for a course.
func (client *DatabaseClient) ListAssignments(span *sentry.Span, guildID string) ([]Assignment, error) {
	span = span.StartChild("listAssignments")
	defer span.Finish()

	course, err := gorm.G[Course](client.DB).Where("guild_id = ?", guildID).First(span.Context())
	if err != nil {
		return []Assignment{}, stacktrace.Propagate(err, "error reading course by guild id: %s", guildID)
	}

	return gorm.G[Assignment](client.DB).Where("course_id = ?", course.ID).Find(span.Context())
}

// CreateAssignment creates a new assignment in the backend.
func (client *DatabaseClient) CreateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {

	err := gorm.G[Assignment](client.DB).Create(span.Context(), &assignment)
	if err != nil {
		return Assignment{}, stacktrace.Propagate(err, "error creating assignment in course: %s", assignment.Course.GuildID)
	}

	return gorm.G[Assignment](client.DB).Where("id = ?", assignment.ID).First(span.Context())
}

// UpdateAssignment updates an existing assignment in the backend.
func (client *DatabaseClient) UpdateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {
	span = span.StartChild("updateAssignment")
	defer span.Finish()

	_, err := gorm.G[Assignment](client.DB).Where("id = ?", assignment.ID).Updates(span.Context(), assignment)
	if err != nil {
		return Assignment{}, stacktrace.Propagate(err, "error updating assignment with id: %d", assignment.ID)
	}

	return gorm.G[Assignment](client.DB).Where("id = ?", assignment.ID).First(span.Context())
}

// DeleteAssignment deletes an assignment from the backend.
func (client *DatabaseClient) DeleteAssignment(span *sentry.Span, assignmentID string) error {
	span = span.StartChild("deleteAssignment")
	defer span.Finish()

	_, err := gorm.G[Assignment](client.DB).Where("id = ?", assignmentID).Delete(span.Context())
	if err != nil {
		return stacktrace.Propagate(err, "error deleting assignment with id: %s", assignmentID)
	}

	return nil
}
