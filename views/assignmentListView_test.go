package views

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/stretchr/testify/suite"
)

// AssignmentListViewTestSuite tests assignment list view functionality
type AssignmentListViewTestSuite struct {
	suite.Suite
	mockMember *discordgo.Member
}

func TestAssignmentListView(t *testing.T) {
	suite.Run(t, new(AssignmentListViewTestSuite))
}

func (s *AssignmentListViewTestSuite) SetupTest() {
	s.mockMember = &discordgo.Member{
		User: &discordgo.User{
			Username: "TestUser",
			ID:       "123456789",
		},
	}
}

// Test AssignmentsListView with empty assignment list
func (s *AssignmentListViewTestSuite) TestAssignmentsListViewEmpty() {
	emptyAssignments := []clients.Assignment{}
	embed := AssignmentsListView(s.mockMember, emptyAssignments)

	s.NotNil(embed)
	s.Equal("assignments", embed.Title)
	s.Equal("0 assignments in course", embed.Description)
	s.Equal("TestUser", embed.Author.Name)
	s.Len(embed.Fields, 0) // No fields for empty list
}

// Test AssignmentsListView with single assignment
func (s *AssignmentListViewTestSuite) TestAssignmentsListViewSingle() {
	testTime := time.Now()
	singleAssignment := []clients.Assignment{
		{
			ID:   "test-id-1",
			Name: "Test Assignment 1",
			Due:  testTime,
			URL:  "https://example.com/1",
		},
	}
	embed := AssignmentsListView(s.mockMember, singleAssignment)

	s.NotNil(embed)
	s.Equal("assignments", embed.Title)
	s.Equal("1 assignments in course", embed.Description)
	s.Len(embed.Fields, 3) // ID, Name, Due fields
	s.Equal("id: test-id-1", embed.Fields[0].Name)
	s.Equal("Test Assignment 1", embed.Fields[0].Value)
	s.Equal("due:", embed.Fields[1].Name)
	s.Equal(testTime.Format(time.RFC1123), embed.Fields[1].Value)
}

// Test AssignmentsListView with multiple assignments
func (s *AssignmentListViewTestSuite) TestAssignmentsListViewMultiple() {
	testTime1 := time.Now()
	testTime2 := testTime1.Add(24 * time.Hour)
	multipleAssignments := []clients.Assignment{
		{
			ID:   "test-id-1",
			Name: "Test Assignment 1",
			Due:  testTime1,
			URL:  "https://example.com/1",
		},
		{
			ID:   "test-id-2",
			Name: "Test Assignment 2",
			Due:  testTime2,
			URL:  "https://example.com/2",
		},
	}
	embed := AssignmentsListView(s.mockMember, multipleAssignments)

	s.NotNil(embed)
	s.Equal("assignments", embed.Title)
	s.Equal("2 assignments in course", embed.Description)
	s.Len(embed.Fields, 6) // 3 fields per assignment
}

// Test AssignmentsListActions
func (s *AssignmentListViewTestSuite) TestAssignmentsListActions() {
	actions := AssignmentsListActions()

	s.NotNil(actions)
	s.Len(actions.Components, 1)

	button := actions.Components[0].(discordgo.Button)
	s.Equal("➕", button.Emoji.Name)
	s.Equal("add", button.Label)
	s.Equal(discordgo.PrimaryButton, button.Style)
	s.Equal("addAssignmentAction", button.CustomID)
}
