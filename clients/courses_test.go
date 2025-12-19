package clients_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dragonejt/hakase-discord/clients"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/suite"
)

// CourseTestSuite tests the actual BackendClient implementation for courses
type CourseTestSuite struct {
	suite.Suite
	testServer    *httptest.Server
	backendClient *clients.BackendClient
	testSpan      *sentry.Span
}

func TestCourses(t *testing.T) {
	suite.Run(t, new(CourseTestSuite))
}

func (testSuite *CourseTestSuite) SetupTest() {
	// Create a test HTTP server that we can control
	testSuite.testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will be overridden in individual tests
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))

	testSuite.backendClient = &clients.BackendClient{
		URL:        testSuite.testServer.URL,
		AuthToken:  "test-token",
		HTTPClient: testSuite.testServer.Client(),
	}
	testSuite.testSpan = sentry.StartSpan(context.Background(), testSuite.T().Name())
}

func (testSuite *CourseTestSuite) TearDownTest() {
	testSuite.testSpan.Finish()
	testSuite.testServer.Close()
}

// Test ReadCourse with successful API response
func (testSuite *CourseTestSuite) TestReadCourseSuccess() {
	testCourse := clients.Course{
		Course:        "course-1",
		CourseID:      "course-1",
		NotifyChannel: "channel-1",
		NotifyGroup:   "group-1",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		testSuite.Equal(http.MethodGet, r.Method)
		testSuite.Equal("/objects/Course/course-1", r.URL.Path)
		testSuite.Equal("Bearer test-token", r.Header.Get("authorization"))

		// Return mock response
		w.WriteHeader(http.StatusOK)
		responseJSON := fmt.Sprintf(`{
			"Course": "%s",
			"courseId": "%s",
			"notifyChannel": "%s",
			"notifyGroup": "%s"
		}`, testCourse.Course, testCourse.CourseID, testCourse.NotifyChannel, testCourse.NotifyGroup)
		w.Write([]byte(responseJSON))
	})

	result, err := testSuite.backendClient.ReadCourse(testSuite.testSpan, "course-1")
	testSuite.NoError(err)
	testSuite.Equal(testCourse.Course, result.Course)
	testSuite.Equal(testCourse.CourseID, result.CourseID)
	testSuite.Equal(testCourse.NotifyChannel, result.NotifyChannel)
	testSuite.Equal(testCourse.NotifyGroup, result.NotifyGroup)
}

// Test ReadCourse with error response
func (testSuite *CourseTestSuite) TestReadCourseError() {
	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	})

	result, err := testSuite.backendClient.ReadCourse(testSuite.testSpan, "invalid-course")
	testSuite.Error(err)
	testSuite.Equal("", result.CourseID)
}

// Test CreateCourse with successful API response
func (testSuite *CourseTestSuite) TestCreateCourseSuccess() {
	testCourse := clients.Course{
		Course:        "New Course",
		CourseID:      "course-1",
		NotifyChannel: "channel-1",
		NotifyGroup:   "group-1",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		testSuite.Equal(http.MethodPost, r.Method)
		testSuite.Equal("/actions/create-course/apply", r.URL.Path)
		testSuite.Equal("Bearer test-token", r.Header.Get("authorization"))

		// Return mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"validation": {
				"result": "VALID"
			}
		}`))
	})

	err := testSuite.backendClient.CreateCourse(testSuite.testSpan, testCourse)
	testSuite.NoError(err)
}

// Test CreateCourse with error response
func (testSuite *CourseTestSuite) TestCreateCourseError() {
	testCourse := clients.Course{
		Course:        "Invalid Course",
		CourseID:      "invalid-course",
		NotifyChannel: "channel-1",
		NotifyGroup:   "group-1",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"validation": {
				"result": "INVALID"
			}
		}`))
	})

	err := testSuite.backendClient.CreateCourse(testSuite.testSpan, testCourse)
	testSuite.Error(err)
}

// Test UpdateCourse with successful API response
func (testSuite *CourseTestSuite) TestUpdateCourseSuccess() {
	testCourse := clients.Course{
		Course:        "Updated Course",
		CourseID:      "existing-course-id",
		NotifyChannel: "channel-1",
		NotifyGroup:   "group-1",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		testSuite.Equal(http.MethodPost, r.Method)
		testSuite.Equal("/actions/edit-course/apply", r.URL.Path)
		testSuite.Equal("Bearer test-token", r.Header.Get("authorization"))

		// Return mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"validation": {
				"result": "VALID"
			}
		}`))
	})

	err := testSuite.backendClient.UpdateCourse(testSuite.testSpan, testCourse)
	testSuite.NoError(err)
}

// Test UpdateCourse with error response
func (testSuite *CourseTestSuite) TestUpdateCourseError() {
	testCourse := clients.Course{
		Course:        "Invalid Update",
		CourseID:      "invalid-course-id",
		NotifyChannel: "channel-1",
		NotifyGroup:   "group-1",
	}

	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	})

	err := testSuite.backendClient.UpdateCourse(testSuite.testSpan, testCourse)
	testSuite.Error(err)
}

// Test DeleteCourse with successful API response
func (testSuite *CourseTestSuite) TestDeleteCourseSuccess() {
	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		testSuite.Equal(http.MethodPost, r.Method)
		testSuite.Equal("/actions/delete-course/apply", r.URL.Path)
		testSuite.Equal("Token test-token", r.Header.Get("authorization"))

		// Return mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"validation": {
				"result": "VALID"
			}
		}`))
	})

	err := testSuite.backendClient.DeleteCourse(testSuite.testSpan, "course-1")
	testSuite.NoError(err)
}

// Test DeleteCourse with error response
func (testSuite *CourseTestSuite) TestDeleteCourseError() {
	// Reset the test server handler for this specific test
	testSuite.testServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	})

	err := testSuite.backendClient.DeleteCourse(testSuite.testSpan, "invalid-course")
	testSuite.Error(err)
}
