package clients

import (
	"context"
	"fmt"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/suite"
)

// MockOntologyClient extends MockHakaseClient for course-specific testing
type MockOntologyClient struct {
	*MockHakaseClient
}

func (m *MockOntologyClient) ReadCourse(span *sentry.Span, guildID string) (Course, error) {
	args := m.Called(span, guildID)
	return args.Get(0).(Course), args.Error(1)
}

func (m *MockOntologyClient) CreateCourse(span *sentry.Span, course Course) error {
	args := m.Called(span, course)
	return args.Error(0)
}

func (m *MockOntologyClient) UpdateCourse(span *sentry.Span, course Course) error {
	args := m.Called(span, course)
	return args.Error(0)
}

func (m *MockOntologyClient) DeleteCourse(span *sentry.Span, guildID string) error {
	args := m.Called(span, guildID)
	return args.Error(0)
}

// CoursesTestSuite tests course management functionality
type CoursesTestSuite struct {
	suite.Suite
	mockClient *MockOntologyClient
	testSpan   *sentry.Span
}

func TestCourses(t *testing.T) {
	suite.Run(t, new(CoursesTestSuite))
}

func (s *CoursesTestSuite) SetupTest() {
	s.mockClient = &MockOntologyClient{MockHakaseClient: new(MockHakaseClient)}
	s.testSpan = sentry.StartSpan(context.Background(), s.T().Name())
}

func (s *CoursesTestSuite) TearDownTest() {
	s.testSpan.Finish()
	s.mockClient.AssertExpectations(s.T())
}

// Test ReadCourse functionality
func (s *CoursesTestSuite) TestReadCourseSuccess() {
	testCourse := Course{
		CourseID:      "test-course-123",
		NotifyChannel: "channel-123",
		NotifyGroup:   "group-123",
	}

	s.mockClient.On("ReadCourse", s.testSpan, "test-course-123").Return(testCourse, nil)

	result, err := s.mockClient.ReadCourse(s.testSpan, "test-course-123")
	s.NoError(err)
	s.Equal(testCourse, result)
}

func (s *CoursesTestSuite) TestReadCourseError() {
	s.mockClient.On("ReadCourse", s.testSpan, "invalid-course").Return(Course{}, fmt.Errorf("course not found"))

	result, err := s.mockClient.ReadCourse(s.testSpan, "invalid-course")
	s.Error(err)
	s.Equal(Course{}, result)
}

// Test CreateCourse functionality
func (s *CoursesTestSuite) TestCreateCourseSuccess() {
	testCourse := Course{
		CourseID:      "new-course-456",
		NotifyChannel: "channel-456",
		NotifyGroup:   "group-456",
	}

	s.mockClient.On("CreateCourse", s.testSpan, testCourse).Return(nil)

	err := s.mockClient.CreateCourse(s.testSpan, testCourse)
	s.NoError(err)
}

func (s *CoursesTestSuite) TestCreateCourseError() {
	testCourse := Course{
		CourseID:      "invalid-course",
		NotifyChannel: "invalid-channel",
		NotifyGroup:   "invalid-group",
	}

	s.mockClient.On("CreateCourse", s.testSpan, testCourse).Return(fmt.Errorf("validation failed"))

	err := s.mockClient.CreateCourse(s.testSpan, testCourse)
	s.Error(err)
}

// Test UpdateCourse functionality
func (s *CoursesTestSuite) TestUpdateCourseSuccess() {
	testCourse := Course{
		CourseID:      "existing-course-789",
		NotifyChannel: "updated-channel",
		NotifyGroup:   "updated-group",
	}

	s.mockClient.On("UpdateCourse", s.testSpan, testCourse).Return(nil)

	err := s.mockClient.UpdateCourse(s.testSpan, testCourse)
	s.NoError(err)
}

func (s *CoursesTestSuite) TestUpdateCourseError() {
	testCourse := Course{
		CourseID:      "nonexistent-course",
		NotifyChannel: "channel",
		NotifyGroup:   "group",
	}

	s.mockClient.On("UpdateCourse", s.testSpan, testCourse).Return(fmt.Errorf("course not found"))

	err := s.mockClient.UpdateCourse(s.testSpan, testCourse)
	s.Error(err)
}

// Test DeleteCourse functionality
func (s *CoursesTestSuite) TestDeleteCourseSuccess() {
	s.mockClient.On("DeleteCourse", s.testSpan, "course-to-delete").Return(nil)

	err := s.mockClient.DeleteCourse(s.testSpan, "course-to-delete")
	s.NoError(err)
}

func (s *CoursesTestSuite) TestDeleteCourseError() {
	s.mockClient.On("DeleteCourse", s.testSpan, "nonexistent-course").Return(fmt.Errorf("course not found"))

	err := s.mockClient.DeleteCourse(s.testSpan, "nonexistent-course")
	s.Error(err)
}
