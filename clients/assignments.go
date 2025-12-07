// Package clients implements backend API operations for assignments.
package clients

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

type Assignment struct {
	ID       string
	CourseID string
	Name     string
	Due      time.Time
	URL      string
}

// ReadAssignment retrieves an assignment by its ID from the backend.
func (client *OntologyClient) ReadAssignment(span *sentry.Span, assignmentID string) (Assignment, error) {
	span = span.StartChild("readAssignment")
	defer span.Finish()

	response, err := client.AssignmentGetAssignmentWithResponse(span.Context(), assignmentID, new(AssignmentGetAssignmentParams))
	if err != nil || response.JSON200 == nil {
		return Assignment{}, stacktrace.Propagate(err, "failed to read assignment with assignment ID: %s", assignmentID)
	}

	return Assignment{
		ID:       *response.JSON200.Id,
		CourseID: *response.JSON200.CourseId,
		Name:     *response.JSON200.Name,
		Due:      *response.JSON200.Due,
		URL:      *response.JSON200.Url,
	}, nil
}

// ListAssignments lists all assignments for a course.
func (client *OntologyClient) ListAssignments(span *sentry.Span, guildID string) ([]Assignment, error) {
	span = span.StartChild("listAssignments")
	defer span.Finish()

	assignments := []Assignment{}
	response, err := client.CourseListCourseLinkedAssignmentsWithResponse(span.Context(), guildID, new(CourseListCourseLinkedAssignmentsParams))
	if err != nil || response.JSON200 == nil {
		return assignments, stacktrace.Propagate(err, "failed to list assignments with course ID: %s", guildID)
	}

	for _, assignment := range *response.JSON200.Data {
		assignments = append(assignments, Assignment{
			ID:       *assignment.Id,
			CourseID: *assignment.CourseId,
			Name:     *assignment.Name,
			Due:      *assignment.Due,
			URL:      *assignment.Url,
		})
	}

	return assignments, nil
}

// CreateAssignment creates a new assignment in the backend.
func (client *OntologyClient) CreateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {
	span = span.StartChild("createAssignment")
	defer span.Finish()

	mode := VALIDATEANDEXECUTE
	returnEdits := ALL
	createResponse, err := client.CreateAssignmentApplyCreateAssignmentWithResponse(span.Context(), new(CreateAssignmentApplyCreateAssignmentParams), CreateAssignmentApplyCreateAssignmentJSONRequestBody{
		Parameters: OsdkCreateAssignmentParameters{
			CourseId: assignment.CourseID,
			Due:      assignment.Due,
			Name:     assignment.Name,
			Url:      assignment.URL,
		},
		Options: OntologiesApplyActionRequestOptions{
			Mode:        &mode,
			ReturnEdits: &returnEdits,
		},
	})
	if err != nil || createResponse.JSON200 == nil {
		return Assignment{}, stacktrace.Propagate(err, "failed to create course")
	}
	if createResponse.JSON200.Validation.Result == INVALID {
		return Assignment{}, stacktrace.NewError("failed validation: %s", fmt.Sprint(createResponse.JSON200.Validation.Parameters))
	}

	edits := []OntologyEdit{}
	objectEdits, err := createResponse.JSON200.Edits.AsOntologiesObjectEdits()
	if err != nil {
		return Assignment{}, stacktrace.Propagate(err, "failed to get ontologies object edits")
	}
	for _, objectEdit := range *objectEdits.Edits {
		edit := OntologyEdit{}
		err := json.Unmarshal(objectEdit.union, &edit)
		if err != nil {
			return Assignment{}, stacktrace.Propagate(err, "failed to unmarshal to OntologyEdit")
		}
		edits = append(edits, edit)
	}

	assignmentID := edits[0].PrimaryKey
	readResponse, err := client.AssignmentGetAssignmentWithResponse(span.Context(), assignmentID, new(AssignmentGetAssignmentParams))
	if err != nil || readResponse.JSON200 == nil {
		return Assignment{}, stacktrace.Propagate(err, "failed to read assignment with assignment ID: %s", assignmentID)
	}
	return Assignment{
		ID:       *readResponse.JSON200.Id,
		CourseID: *readResponse.JSON200.CourseId,
		Name:     *readResponse.JSON200.Name,
		Due:      *readResponse.JSON200.Due,
		URL:      *readResponse.JSON200.Url,
	}, nil
}

// UpdateAssignment updates an existing assignment in the backend.
func (client *OntologyClient) UpdateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {
	span = span.StartChild("updateAssignment")
	defer span.Finish()

	mode := VALIDATEANDEXECUTE
	returnEdits := ALL
	editResponse, err := client.EditAssignmentApplyEditAssignmentWithResponse(span.Context(), new(EditAssignmentApplyEditAssignmentParams), EditAssignmentApplyEditAssignmentJSONRequestBody{
		Parameters: OsdkEditAssignmentParameters{
			Due:  assignment.Due,
			Name: assignment.Name,
			Url:  assignment.URL,
		},
		Options: OntologiesApplyActionRequestOptions{
			Mode:        &mode,
			ReturnEdits: &returnEdits,
		},
	})
	if err != nil || editResponse.JSON200 == nil {
		return Assignment{}, stacktrace.Propagate(err, "failed to update assignment with assignment ID: %s", assignment.ID)
	}
	if editResponse.JSON200.Validation.Result == INVALID {
		return Assignment{}, stacktrace.NewError("failed validation: %s", fmt.Sprint(editResponse.JSON200.Validation.Parameters))
	}

	edits := []OntologyEdit{}
	objectEdits, err := editResponse.JSON200.Edits.AsOntologiesObjectEdits()
	if err != nil {
		return Assignment{}, stacktrace.Propagate(err, "failed to get ontologies object edits")
	}
	for _, objectEdit := range *objectEdits.Edits {
		edit := OntologyEdit{}
		err := json.Unmarshal(objectEdit.union, &edit)
		if err != nil {
			return Assignment{}, stacktrace.Propagate(err, "failed to unmarshal to OntologyEdit")
		}
		edits = append(edits, edit)
	}

	assignmentID := edits[0].PrimaryKey
	readResponse, err := client.AssignmentGetAssignmentWithResponse(span.Context(), assignmentID, new(AssignmentGetAssignmentParams))
	if err != nil || readResponse.JSON200 == nil {
		return Assignment{}, stacktrace.Propagate(err, "failed to read assignment with assignment ID: %s", assignmentID)
	}
	return Assignment{
		ID:       *readResponse.JSON200.Id,
		CourseID: *readResponse.JSON200.CourseId,
		Name:     *readResponse.JSON200.Name,
		Due:      *readResponse.JSON200.Due,
		URL:      *readResponse.JSON200.Url,
	}, nil
}

// DeleteAssignment deletes an assignment from the backend.
func (client *OntologyClient) DeleteAssignment(span *sentry.Span, assignmentID string) error {
	span = span.StartChild("deleteAssignment")
	defer span.Finish()

	mode := VALIDATEANDEXECUTE
	returnEdits := ALL
	response, err := client.DeleteAssignmentApplyDeleteAssignmentWithResponse(span.Context(), new(DeleteAssignmentApplyDeleteAssignmentParams), DeleteAssignmentApplyDeleteAssignmentJSONRequestBody{
		Parameters: OsdkDeleteAssignmentParameters{
			Assignment: assignmentID,
		},
		Options: OntologiesApplyActionRequestOptions{
			Mode:        &mode,
			ReturnEdits: &returnEdits,
		},
	})
	if err != nil || response.JSON200 == nil {
		return stacktrace.Propagate(err, "failed to delete assignment with assignment ID: %s", assignmentID)
	}
	if response.JSON200.Validation.Result == INVALID {
		return stacktrace.NewError("failed validation: %s", fmt.Sprint(response.JSON200.Validation.Parameters))
	}

	return nil
}
