// Package clients implements backend API operations for courses.
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

type Course struct {
	Course        string `json:"Course,omitempty"`
	CourseID      string `json:"courseId,omitempty"`
	NotifyChannel string `json:"notifyChannel,omitempty"`
	NotifyGroup   string `json:"notifyGroup,omitempty"`
}

type CourseRequest struct {
	Parameters Course         `json:"parameters"`
	Options    RequestOptions `json:"options"`
}

// ReadCourse retrieves a course by its ID from the backend.
func (backend *BackendClient) ReadCourse(span *sentry.Span, guildID string) (Course, error) {
	span = span.StartChild("readCourse")
	defer span.Finish()

	course := Course{}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodGet, fmt.Sprintf("%s/objects/Course/%s", backend.URL, guildID), nil)
	if err != nil {
		return course, stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", backend.AuthToken))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTPClient.Do(request)
	if err != nil {
		return course, stacktrace.Propagate(err, "failed to execute API request")
	}
	if response.StatusCode != http.StatusOK {
		return course, stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return course, stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &course)
	if err != nil {
		return course, stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}

	return course, nil
}

// CreateCourse creates a new course in the backend.
func (backend *BackendClient) CreateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("createCourse")
	defer span.Finish()

	courseRequest := CourseRequest{
		Parameters: course,
		Options: RequestOptions{
			ReturnEdits: "ALL",
		},
	}
	courseResponse := ActionResponse{}

	jsonBody, err := json.Marshal(courseRequest)
	if err != nil {
		return stacktrace.Propagate(err, "failed to marshal course")
	}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodPost, fmt.Sprintf("%s/actions/create-course/apply", backend.URL), bytes.NewReader(jsonBody))
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
	if response.StatusCode != http.StatusOK {
		return stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &courseResponse)
	if err != nil {
		return stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}
	if courseResponse.Validation.Result != "VALID" {
		return stacktrace.NewError("request validation result: %s", courseResponse.Validation.Result)
	}

	return nil
}

// UpdateCourse updates an existing course in the backend.
func (backend *BackendClient) UpdateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("updateCourse")
	defer span.Finish()

	course.CourseID = ""
	courseRequest := CourseRequest{
		Parameters: course,
		Options: RequestOptions{
			ReturnEdits: "ALL",
		},
	}
	courseResponse := ActionResponse{}

	jsonBody, err := json.Marshal(courseRequest)
	if err != nil {
		return stacktrace.Propagate(err, "failed to marshal course")
	}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodPost, fmt.Sprintf("%s/actions/edit-course/apply", backend.URL), bytes.NewReader(jsonBody))
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
	if response.StatusCode != http.StatusOK {
		return stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &courseResponse)
	if err != nil {
		return stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}
	if courseResponse.Validation.Result != "VALID" {
		return stacktrace.NewError("request validation result: %s", courseResponse.Validation.Result)
	}

	return nil
}

// DeleteCourse deletes a course from the backend.
func (backend *BackendClient) DeleteCourse(span *sentry.Span, guildID string) error {
	span = span.StartChild("deleteCourse")
	defer span.Finish()

	courseRequest := CourseRequest{
		Parameters: Course{
			Course: guildID,
		},
		Options: RequestOptions{
			ReturnEdits: "ALL",
		},
	}
	courseResponse := ActionResponse{}

	jsonBody, err := json.Marshal(courseRequest)
	if err != nil {
		return stacktrace.Propagate(err, "failed to marshal course")
	}

	request, err := http.NewRequestWithContext(span.Context(), http.MethodPost, fmt.Sprintf("%s/actions/delete-course/apply", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.AuthToken))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTPClient.Do(request)
	if err != nil {
		return stacktrace.Propagate(err, "failed to execute API request")
	}
	if response.StatusCode != http.StatusOK {
		return stacktrace.NewError("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &courseResponse)
	if err != nil {
		return stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}
	if courseResponse.Validation.Result != "VALID" {
		return stacktrace.NewError("request validation result: %s", courseResponse.Validation.Result)
	}

	return nil
}
