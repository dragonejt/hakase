// Package clients implements backend API operations for assignments.
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

type Assignment struct {
	Assignment string    `json:"Assignment,omitempty"`
	ID         string    `json:"id,omitempty"`
	CourseID   string    `json:"courseId,omitempty"`
	Status     string    `json:"status,omitempty"`
	Name       string    `json:"name,omitempty"`
	Due        time.Time `json:"due,omitzero"`
	URL        string    `json:"url,omitempty"`
}

var AssignmentStatus = []string{"not due", "one day", "one hour"}

type AssignmentRequest struct {
	Parameters Assignment     `json:"parameters"`
	Options    RequestOptions `json:"options"`
}

type AssignmentsResponse struct {
	Data []Assignment `json:"data,omitempty"`
}

type SearchAssignmentsRequest struct {
	Where SearchAssignmentsQuery `json:"where"`
}

type SearchAssignmentsQuery struct {
	Type  string `json:"type,omitempty"`
	Field string `json:"field,omitempty"`
	Value any    `json:"value,omitempty"`
}

// ReadAssignment retrieves an assignment by its ID from the backend.
func (backend *BackendClient) ReadAssignment(span *sentry.Span, assignmentID string) (Assignment, error) {
	span = span.StartChild("readAssignment")
	defer span.Finish()

	assignment := Assignment{}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodGet, fmt.Sprintf("%s/objects/Assignment/%s", backend.URL, assignmentID), nil)
	if err != nil {
		return assignment, stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", backend.AuthToken))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTPClient.Do(request)
	if err != nil {
		return assignment, stacktrace.Propagate(err, "failed to execute API request")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return assignment, stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return assignment, stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &assignment)
	if err != nil {
		return assignment, stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}

	return assignment, nil
}

// ListAssignments lists all assignments for a course.
func (backend *BackendClient) ListAssignments(span *sentry.Span, guildID string) ([]Assignment, error) {
	span = span.StartChild("listAssignments")
	defer span.Finish()

	listAssignmentsResponse := AssignmentsResponse{}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodGet, fmt.Sprintf("%s/objects/Course/%s/links/assignments", backend.URL, guildID), nil)
	if err != nil {
		return listAssignmentsResponse.Data, stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", backend.AuthToken))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTPClient.Do(request)
	if err != nil {
		return listAssignmentsResponse.Data, stacktrace.Propagate(err, "failed to execute API request")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return listAssignmentsResponse.Data, stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return listAssignmentsResponse.Data, stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &listAssignmentsResponse)
	if err != nil {
		return listAssignmentsResponse.Data, stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}

	return listAssignmentsResponse.Data, nil
}

// SearchAssignments searches through assignments based on a filter.
func (backend *BackendClient) SearchAssignments(span *sentry.Span, searchQuery SearchAssignmentsQuery) ([]Assignment, error) {
	span = span.StartChild("searchAssignments")
	defer span.Finish()

	searchAssignmentsRequest := SearchAssignmentsRequest{
		Where: searchQuery,
	}
	searchAssignmentsResponse := AssignmentsResponse{}

	jsonBody, err := json.Marshal(searchAssignmentsRequest)
	if err != nil {
		return searchAssignmentsResponse.Data, stacktrace.Propagate(err, "failed to marshal assignment")
	}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodPost, fmt.Sprintf("%s/objects/Assignment/search", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return searchAssignmentsResponse.Data, stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", backend.AuthToken))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTPClient.Do(request)
	if err != nil {
		return searchAssignmentsResponse.Data, stacktrace.Propagate(err, "failed to execute API request")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return searchAssignmentsResponse.Data, stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return searchAssignmentsResponse.Data, stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &searchAssignmentsResponse)
	if err != nil {
		return searchAssignmentsResponse.Data, stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}

	return searchAssignmentsResponse.Data, nil
}

// CreateAssignment creates a new assignment in the backend.
func (backend *BackendClient) CreateAssignment(span *sentry.Span, assignment Assignment) (string, error) {
	span = span.StartChild("createAssignment")
	defer span.Finish()

	assignment.Status = AssignmentStatus[0]
	assignmentRequest := AssignmentRequest{
		Parameters: assignment,
		Options: RequestOptions{
			ReturnEdits: "ALL",
		},
	}
	assignmentResponse := ActionResponse{}

	jsonBody, err := json.Marshal(assignmentRequest)
	if err != nil {
		return "", stacktrace.Propagate(err, "failed to marshal assignment")
	}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodPost, fmt.Sprintf("%s/actions/create-assignment/apply", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return "", stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", backend.AuthToken))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTPClient.Do(request)
	if err != nil {
		return "", stacktrace.Propagate(err, "failed to execute API request")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &assignmentResponse)
	if err != nil {
		return "", stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}
	if assignmentResponse.Validation.Result != "VALID" {
		return "", stacktrace.NewError("request validation result: %s", assignmentResponse.Validation.Result)
	}

	return assignmentResponse.Edits.Edits[0].PrimaryKey, nil
}

// UpdateAssignment updates an existing assignment in the backend.
func (backend *BackendClient) UpdateAssignment(span *sentry.Span, assignment Assignment) error {
	span = span.StartChild("updateAssignment")
	defer span.Finish()

	assignment.ID = ""
	assignmentRequest := AssignmentRequest{
		Parameters: assignment,
		Options: RequestOptions{
			ReturnEdits: "ALL",
		},
	}
	assignmentResponse := ActionResponse{}

	jsonBody, err := json.Marshal(assignmentRequest)
	if err != nil {
		return stacktrace.Propagate(err, "failed to marshal assignment")
	}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodPost, fmt.Sprintf("%s/actions/edit-assignment/apply", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", backend.AuthToken))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTPClient.Do(request)
	if err != nil {
		return stacktrace.Propagate(err, "failed to execute API request")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &assignmentResponse)
	if err != nil {
		return stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}
	if assignmentResponse.Validation.Result != "VALID" {
		return stacktrace.NewError("request validation result: %s", assignmentResponse.Validation.Result)
	}

	return nil
}

// DeleteAssignment deletes an assignment from the backend.
func (backend *BackendClient) DeleteAssignment(span *sentry.Span, assignmentID string) error {
	span = span.StartChild("deleteAssignment")
	defer span.Finish()

	assignmentRequest := AssignmentRequest{
		Parameters: Assignment{
			Assignment: assignmentID,
		},
		Options: RequestOptions{
			ReturnEdits: "ALL",
		},
	}
	assignmentResponse := ActionResponse{}

	jsonBody, err := json.Marshal(assignmentRequest)
	if err != nil {
		return stacktrace.Propagate(err, "failed to marshal assignment")
	}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodPost, fmt.Sprintf("%s/actions/delete-assignment/apply", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", backend.AuthToken))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTPClient.Do(request)
	if err != nil {
		return stacktrace.Propagate(err, "failed to execute API request")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &assignmentResponse)
	if err != nil {
		return stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}
	if assignmentResponse.Validation.Result != "VALID" {
		return stacktrace.NewError("request validation result: %s", assignmentResponse.Validation.Result)
	}

	return nil
}
