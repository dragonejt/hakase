package clients_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dragonejt/hakase-discord/clients"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/suite"
)

// AssignmentTestSuite tests the actual BackendClient implementation for assignments
type AssignmentTestSuite struct {
	suite.Suite
	testServer    *httptest.Server
	backendClient *clients.BackendClient
	testSpan      *sentry.Span
}

func TestAssignments(t *testing.T) {
	suite.Run(t, new(AssignmentTestSuite))
}

func (testSuite *AssignmentTestSuite) SetupTest() {
	// Create a test HTTP server that we can control
	testSuite.testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will be overridden in individual tests
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}))

	testSuite.backendClient = &clients.BackendClient{
		URL:        testSuite.testServer.URL,
		AuthToken:  "test-token",
		HTTPClient: testSuite.testServer.Client(),
	}
	testSuite.testSpan = sentry.StartSpan(context.Background(), testSuite.T().Name())
}

func (testSuite *AssignmentTestSuite) TearDownTest() {
	testSuite.testSpan.Finish()
	testSuite.testServer.Close()
}

// Test ReadAssignment with successful API response
func (testSuite *AssignmentTestSuite) TestReadAssignmentSuccess() {
	testTime := time.Now()
	testAssignment := clients.Assignment{
		ID:       "test-id",
		CourseID: "test-course",
		Name:     "Test Assignment",
		Due:      &testTime,
		URL:      "https://example.com",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		testSuite.Equal(http.MethodGet, r.Method)
		testSuite.Equal("/objects/Assignment/test-id", r.URL.Path)
		testSuite.Equal("Bearer test-token", r.Header.Get("authorization"))

		// Return mock response
		w.WriteHeader(http.StatusOK)
		responseJSON := fmt.Sprintf(`{
			"id": "%s",
			"courseId": "%s",
			"name": "%s",
			"due": "%s",
			"url": "%s"
		}`, testAssignment.ID, testAssignment.CourseID, testAssignment.Name, testTime.Format(time.RFC3339), testAssignment.URL)
		if _, err := w.Write([]byte(responseJSON)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	result, err := testSuite.backendClient.ReadAssignment(testSuite.testSpan, "test-id")
	testSuite.NoError(err)
	testSuite.Equal(testAssignment.ID, result.ID)
	testSuite.Equal(testAssignment.Name, result.Name)
	testSuite.Equal(testAssignment.CourseID, result.CourseID)
}

// Test ReadAssignment with error response
func (testSuite *AssignmentTestSuite) TestReadAssignmentError() {
	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`{"error": "not found"}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	result, err := testSuite.backendClient.ReadAssignment(testSuite.testSpan, "invalid-id")
	testSuite.Error(err)
	testSuite.Equal("", result.ID)
}

// Test ListAssignments with successful API response
func (testSuite *AssignmentTestSuite) TestListAssignmentsSuccess() {
	testTime1 := time.Now()
	testTime2 := testTime1.Add(24 * time.Hour)
	testAssignments := []clients.Assignment{
		{
			ID:       "assignment-1",
			CourseID: "course-1",
			Name:     "Assignment 1",
			Due:      &testTime1,
			URL:      "https://example.com/1",
		},
		{
			ID:       "assignment-2",
			CourseID: "course-1",
			Name:     "Assignment 2",
			Due:      &testTime2,
			URL:      "https://example.com/2",
		},
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		testSuite.Equal(http.MethodGet, r.Method)
		testSuite.Equal("/objects/Course/course-1/links/assignments", r.URL.Path)
		testSuite.Equal("Bearer test-token", r.Header.Get("authorization"))

		// Return mock response
		w.WriteHeader(http.StatusOK)
		responseJSON := fmt.Sprintf(`{
			"data": [
				{
					"id": "%s",
					"courseId": "%s",
					"name": "%s",
					"due": "%s",
					"url": "%s"
				},
				{
					"id": "%s",
					"courseId": "%s",
					"name": "%s",
					"due": "%s",
					"url": "%s"
				}
			]
		}`, testAssignments[0].ID, testAssignments[0].CourseID, testAssignments[0].Name, testTime1.Format(time.RFC3339), testAssignments[0].URL,
			testAssignments[1].ID, testAssignments[1].CourseID, testAssignments[1].Name, testTime2.Format(time.RFC3339), testAssignments[1].URL)
		if _, err := w.Write([]byte(responseJSON)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	result, err := testSuite.backendClient.ListAssignments(testSuite.testSpan, "course-1")
	testSuite.NoError(err)
	testSuite.Len(result, 2)
	testSuite.Equal(testAssignments[0].ID, result[0].ID)
	testSuite.Equal(testAssignments[1].ID, result[1].ID)
}

// Test ListAssignments with error response
func (testSuite *AssignmentTestSuite) TestListAssignmentsError() {
	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`{"error": "not found"}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	result, err := testSuite.backendClient.ListAssignments(testSuite.testSpan, "invalid-course")
	testSuite.Error(err)
	testSuite.Equal(0, len(result))
}

// Test CreateAssignment with successful API response
func (testSuite *AssignmentTestSuite) TestCreateAssignmentSuccess() {
	testTime := time.Now()
	testAssignment := clients.Assignment{
		CourseID: "course-1",
		Name:     "New Assignment",
		Due:      &testTime,
		URL:      "https://example.com/new",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		testSuite.Equal(http.MethodPost, r.Method)
		testSuite.Equal("/actions/create-assignment/apply", r.URL.Path)
		testSuite.Equal("Bearer test-token", r.Header.Get("authorization"))
		testSuite.Equal("application/json", r.Header.Get("content-type"))

		// Return mock response
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{
			"validation": {
				"result": "VALID"
			},
			"edits": {
				"edits": [
					{
						"primaryKey": "created-assignment-id"
					}
				]
			}
		}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	createdID, err := testSuite.backendClient.CreateAssignment(testSuite.testSpan, testAssignment)
	testSuite.NoError(err)
	testSuite.Equal("created-assignment-id", createdID)
}

// Test CreateAssignment with error response
func (testSuite *AssignmentTestSuite) TestCreateAssignmentError() {
	testTime := time.Now()
	testAssignment := clients.Assignment{
		CourseID: "invalid-course",
		Name:     "Invalid Assignment",
		Due:      &testTime,
		URL:      "https://example.com/invalid",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{
			"validation": {
				"result": "INVALID"
			}
		}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	createdID, err := testSuite.backendClient.CreateAssignment(testSuite.testSpan, testAssignment)
	testSuite.Error(err)
	testSuite.Equal("", createdID)
}

// Test UpdateAssignment with successful API response
func (testSuite *AssignmentTestSuite) TestUpdateAssignmentSuccess() {
	testTime := time.Now()
	testAssignment := clients.Assignment{
		ID:       "existing-id",
		CourseID: "course-1",
		Name:     "Updated Assignment",
		Due:      &testTime,
		URL:      "https://example.com/updated",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		testSuite.Equal(http.MethodPost, r.Method)
		testSuite.Equal("/actions/edit-assignment/apply", r.URL.Path)
		testSuite.Equal("Bearer test-token", r.Header.Get("authorization"))

		// Return mock response
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{
			"validation": {
				"result": "VALID"
			}
		}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	err := testSuite.backendClient.UpdateAssignment(testSuite.testSpan, testAssignment)
	testSuite.NoError(err)
}

// Test UpdateAssignment with error response
func (testSuite *AssignmentTestSuite) TestUpdateAssignmentError() {
	testTime := time.Now()
	testAssignment := clients.Assignment{
		ID:       "invalid-id",
		CourseID: "course-1",
		Name:     "Invalid Update",
		Due:      &testTime,
		URL:      "https://example.com/invalid",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`{"error": "not found"}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	err := testSuite.backendClient.UpdateAssignment(testSuite.testSpan, testAssignment)
	testSuite.Error(err)
}

// Test DeleteAssignment with successful API response
func (testSuite *AssignmentTestSuite) TestDeleteAssignmentSuccess() {
	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		testSuite.Equal(http.MethodPost, r.Method)
		testSuite.Equal("/actions/delete-assignment/apply", r.URL.Path)
		testSuite.Equal("Bearer test-token", r.Header.Get("authorization"))

		// Return mock response
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{
			"validation": {
				"result": "VALID"
			}
		}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	err := testSuite.backendClient.DeleteAssignment(testSuite.testSpan, "test-assignment-id")
	testSuite.NoError(err)
}

// Test DeleteAssignment with error response
func (testSuite *AssignmentTestSuite) TestDeleteAssignmentError() {
	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`{"error": "not found"}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	err := testSuite.backendClient.DeleteAssignment(testSuite.testSpan, "invalid-assignment-id")
	testSuite.Error(err)
}
