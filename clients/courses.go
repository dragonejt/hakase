// Package clients implements backend API operations for courses.
package clients

import (
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	GuildID       string `gorm:"uniqueIndex"`
	NotifyChannel string
	NotifyGroup   string
}

// ReadCourse retrieves a course by its ID from the backend.
func (client *DatabaseClient) ReadCourse(span *sentry.Span, guildID string) (Course, error) {
	span = span.StartChild("readCourse")
	defer span.Finish()

	return gorm.G[Course](client.DB).Where("guild_id = ?", guildID).First(span.Context())
}

// CreateCourse creates a new course in the backend.
func (client *DatabaseClient) CreateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("createCourse")
	defer span.Finish()

	return gorm.G[Course](client.DB).Create(span.Context(), &course)
}

// UpdateCourse updates an existing course in the backend.
func (client *DatabaseClient) UpdateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("updateCourse")
	defer span.Finish()

	_, err := gorm.G[Course](client.DB).Where("guild_id = ?", course.GuildID).Updates(span.Context(), course)
	if err != nil {
		return stacktrace.Propagate(err, "error updating course with guild id: %s", course.GuildID)
	}
	return nil
}

// DeleteCourse deletes a course from the backend.
func (client *DatabaseClient) DeleteCourse(span *sentry.Span, guildID string) error {
	span = span.StartChild("deleteCourse")
	defer span.Finish()

	_, err := gorm.G[Course](client.DB).Where("guild_id = ?", guildID).Delete(span.Context())
	if err != nil {
		return stacktrace.Propagate(err, "error deleting course with guild id: %s", guildID)
	}

	return nil
}
