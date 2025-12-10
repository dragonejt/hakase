package clients

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockHakaseClient struct {
	HakaseClient
	mock.Mock
}

func (m *MockHakaseClient) ReadAssignment(span *sentry.Span, assignmentID string) (Assignment, error) {
	args := m.Called(span, assignmentID)
	return args.Get(0).(Assignment), args.Error(1)
}

func (m *MockHakaseClient) ListAssignments(span *sentry.Span, courseID string) ([]Assignment, error) {
	args := m.Called(span, courseID)
	return args.Get(0).([]Assignment), args.Error(1)
}

func (m *MockHakaseClient) CreateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {
	args := m.Called(span, assignment)
	return args.Get(0).(Assignment), args.Error(1)
}

func (m *MockHakaseClient) UpdateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {
	args := m.Called(span, assignment)
	return args.Get(0).(Assignment), args.Error(1)
}

func (m *MockHakaseClient) DeleteAssignment(span *sentry.Span, assignmentID string) error {
	args := m.Called(span, assignmentID)
	return args.Error(0)
}

type AssignmentsTestSuite struct {
	suite.Suite
	mockClient *MockHakaseClient
	testSpan   *sentry.Span
}

func TestAssignments(t *testing.T) {
	suite.Run(t, new(AssignmentsTestSuite))
}

func (s *AssignmentsTestSuite) SetupTest() {
	s.mockClient = new(MockHakaseClient)
	s.testSpan = sentry.StartSpan(context.Background(), s.T().Name())
}

func (s *AssignmentsTestSuite) TearDownTest() {
	s.testSpan.Finish()
	s.mockClient.AssertExpectations(s.T())
}

func (s *AssignmentsTestSuite) TestReadAssignmentSuccess() {
	testAssignment := Assignment{
		ID:       "test-id",
		CourseID: "test-course",
		Name:     "Test Assignment",
		Due:      time.Now(),
		URL:      "https://example.com",
	}

	s.mockClient.On("ReadAssignment", s.testSpan, "test-id").Return(testAssignment, nil)

	result, err := s.mockClient.ReadAssignment(s.testSpan, "test-id")
	s.NoError(err)
	s.Equal(testAssignment, result)
}

func (s *AssignmentsTestSuite) TestReadAssignmentError() {
	s.mockClient.On("ReadAssignment", s.testSpan, "invalid-id").Return(Assignment{}, fmt.Errorf("not found"))

	result, err := s.mockClient.ReadAssignment(s.testSpan, "invalid-id")
	s.Error(err)
	s.Equal(Assignment{}, result)
}

func (s *AssignmentsTestSuite) TestListAssignmentsSuccess() {
	testAssignments := []Assignment{
		{ID: "1", CourseID: "course-1", Name: "Assignment 1", Due: time.Now(), URL: "https://example.com/1"},
		{ID: "2", CourseID: "course-1", Name: "Assignment 2", Due: time.Now(), URL: "https://example.com/2"},
	}

	s.mockClient.On("ListAssignments", s.testSpan, "course-1").Return(testAssignments, nil)

	result, err := s.mockClient.ListAssignments(s.testSpan, "course-1")
	s.NoError(err)
	s.Equal(testAssignments, result)
}

func (s *AssignmentsTestSuite) TestListAssignmentsError() {
	s.mockClient.On("ListAssignments", s.testSpan, "invalid-course").Return([]Assignment{}, fmt.Errorf("not found"))

	result, err := s.mockClient.ListAssignments(s.testSpan, "invalid-course")
	s.Error(err)
	s.Equal([]Assignment{}, result)
}

func (s *AssignmentsTestSuite) TestCreateAssignmentSuccess() {
	testAssignment := Assignment{
		ID:       "new-id",
		CourseID: "course-1",
		Name:     "New Assignment",
		Due:      time.Now(),
		URL:      "https://example.com/new",
	}

	s.mockClient.On("CreateAssignment", s.testSpan, testAssignment).Return(testAssignment, nil)

	result, err := s.mockClient.CreateAssignment(s.testSpan, testAssignment)
	s.NoError(err)
	s.Equal(testAssignment, result)
}

func (s *AssignmentsTestSuite) TestCreateAssignmentError() {
	testAssignment := Assignment{
		CourseID: "invalid-course",
		Name:     "Invalid Assignment",
		Due:      time.Now(),
		URL:      "https://example.com/invalid",
	}

	s.mockClient.On("CreateAssignment", s.testSpan, testAssignment).Return(Assignment{}, fmt.Errorf("validation failed"))

	result, err := s.mockClient.CreateAssignment(s.testSpan, testAssignment)
	s.Error(err)
	s.Equal(Assignment{}, result)
}

func (s *AssignmentsTestSuite) TestUpdateAssignmentSuccess() {
	testAssignment := Assignment{
		ID:       "existing-id",
		CourseID: "course-1",
		Name:     "Updated Assignment",
		Due:      time.Now(),
		URL:      "https://example.com/updated",
	}

	s.mockClient.On("UpdateAssignment", s.testSpan, testAssignment).Return(testAssignment, nil)

	result, err := s.mockClient.UpdateAssignment(s.testSpan, testAssignment)
	s.NoError(err)
	s.Equal(testAssignment, result)
}

func (s *AssignmentsTestSuite) TestUpdateAssignmentError() {
	testAssignment := Assignment{
		ID:       "invalid-id",
		CourseID: "course-1",
		Name:     "Invalid Update",
		Due:      time.Now(),
		URL:      "https://example.com/invalid",
	}

	s.mockClient.On("UpdateAssignment", s.testSpan, testAssignment).Return(Assignment{}, fmt.Errorf("not found"))

	result, err := s.mockClient.UpdateAssignment(s.testSpan, testAssignment)
	s.Error(err)
	s.Equal(Assignment{}, result)
}

func (s *AssignmentsTestSuite) TestDeleteAssignmentSuccess() {
	s.mockClient.On("DeleteAssignment", s.testSpan, "test-id").Return(nil)

	err := s.mockClient.DeleteAssignment(s.testSpan, "test-id")
	s.NoError(err)
}

func (s *AssignmentsTestSuite) TestDeleteAssignmentError() {
	s.mockClient.On("DeleteAssignment", s.testSpan, "invalid-id").Return(fmt.Errorf("not found"))

	err := s.mockClient.DeleteAssignment(s.testSpan, "invalid-id")
	s.Error(err)
}
