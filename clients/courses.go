// Package clients implements backend API operations for courses.
package clients

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

type Course struct {
	CourseID      string
	NotifyChannel string
	NotifyGroup   string
}

// ReadCourse retrieves a course by its ID from the backend.
func (client *OntologyClient) ReadCourse(span *sentry.Span, guildID string) (Course, error) {
	span = span.StartChild("readCourse")
	defer span.Finish()

	response, err := client.CourseGetCourseWithResponse(span.Context(), guildID, new(CourseGetCourseParams))
	if err != nil || response.JSON200 == nil {
		return Course{}, stacktrace.Propagate(err, "failed to read course with course ID: %s", guildID)
	}

	return Course{
		CourseID:      *response.JSON200.CourseId,
		NotifyChannel: *response.JSON200.NotifyChannel,
		NotifyGroup:   *response.JSON200.NotifyGroup,
	}, nil
}

// CreateCourse creates a new course in the backend.
func (client *OntologyClient) CreateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("createCourse")
	defer span.Finish()

	mode := VALIDATEANDEXECUTE
	returnEdits := ALL
	response, err := client.CreateCourseApplyCreateCourseWithResponse(span.Context(), new(CreateCourseApplyCreateCourseParams), CreateCourseApplyCreateCourseJSONRequestBody{
		Parameters: OsdkCreateCourseParameters{
			CourseID:      course.CourseID,
			NotifyChannel: &course.NotifyChannel,
			NotifyGroup:   &course.NotifyGroup,
		},
		Options: OntologiesApplyActionRequestOptions{
			Mode:        &mode,
			ReturnEdits: &returnEdits,
		},
	})
	if err != nil || response.JSON200 == nil {
		return stacktrace.Propagate(err, "failed to create course")
	}
	if response.JSON200.Validation.Result == INVALID {
		return stacktrace.NewError("failed validation: %s", fmt.Sprint(response.JSON200.Validation.Parameters))
	}

	return nil
}

// UpdateCourse updates an existing course in the backend.
func (client *OntologyClient) UpdateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("updateCourse")
	defer span.Finish()

	mode := VALIDATEANDEXECUTE
	returnEdits := ALL
	response, err := client.EditCourseApplyEditCourseWithResponse(span.Context(), new(EditCourseApplyEditCourseParams), EditCourseApplyEditCourseJSONRequestBody{
		Parameters: OsdkEditCourseParameters{
			Course:        course.CourseID,
			NotifyChannel: course.NotifyChannel,
			NotifyGroup:   course.NotifyGroup,
		},
		Options: OntologiesApplyActionRequestOptions{
			Mode:        &mode,
			ReturnEdits: &returnEdits,
		},
	})
	if err != nil || response.JSON200 == nil {
		return stacktrace.Propagate(err, "failed to update course")
	}
	if response.JSON200.Validation.Result == INVALID {
		return stacktrace.NewError("failed validation: %s", fmt.Sprint(response.JSON200.Validation.Parameters))
	}

	return nil
}

// DeleteCourse deletes a course from the backend.
func (client *OntologyClient) DeleteCourse(span *sentry.Span, guildID string) error {
	span = span.StartChild("deleteCourse")
	defer span.Finish()

	mode := VALIDATEANDEXECUTE
	returnEdits := ALL
	response, err := client.DeleteCourseApplyDeleteCourseWithResponse(span.Context(), new(DeleteCourseApplyDeleteCourseParams), DeleteCourseApplyDeleteCourseJSONRequestBody{
		Parameters: OsdkDeleteCourseParameters{
			Course: guildID,
		},
		Options: OntologiesApplyActionRequestOptions{
			Mode:        &mode,
			ReturnEdits: &returnEdits,
		},
	})
	if err != nil || response.JSON200 == nil {
		return stacktrace.Propagate(err, "failed to delete course with course ID: %s", guildID)
	}
	if response.JSON200.Validation.Result == INVALID {
		return stacktrace.NewError("failed validation: %s", fmt.Sprint(response.JSON200.Validation.Parameters))
	}

	return nil
}
