package views_test

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
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

func (testSuite *AssignmentListViewTestSuite) SetupTest() {
	testSuite.mockMember = &discordgo.Member{
		User: &discordgo.User{
			Username: "TestUser",
			ID:       "123456789",
		},
	}
}

// Test AssignmentsListView with empty assignment list
func (testSuite *AssignmentListViewTestSuite) TestAssignmentsListViewEmpty() {
	emptyAssignments := []clients.Assignment{}
	embed := views.AssignmentsListView(testSuite.mockMember, emptyAssignments)

	testSuite.NotNil(embed)
	testSuite.Equal("assignments", embed.Title)
	testSuite.Equal("0 assignments in course", embed.Description)
	testSuite.Equal("TestUser", embed.Author.Name)
	testSuite.Len(embed.Fields, 0) // No fields for empty list
}

// Test AssignmentsListView with single assignment
func (testSuite *AssignmentListViewTestSuite) TestAssignmentsListViewSingle() {
	testTime := time.Now()
	singleAssignment := []clients.Assignment{
		{
			ID:   "test-id-1",
			Name: "Test Assignment 1",
			Due:  testTime,
			URL:  "https://example.com/1",
		},
	}
	embed := views.AssignmentsListView(testSuite.mockMember, singleAssignment)

	testSuite.NotNil(embed)
	testSuite.Equal("assignments", embed.Title)
	testSuite.Equal("1 assignments in course", embed.Description)
	testSuite.Len(embed.Fields, 3) // ID, Name, Due fields
	testSuite.Equal("id: test-id-1", embed.Fields[0].Name)
	testSuite.Equal("Test Assignment 1", embed.Fields[0].Value)
	testSuite.Equal("due:", embed.Fields[1].Name)
	testSuite.Equal(testTime.Format(time.RFC1123), embed.Fields[1].Value)
}

// Test AssignmentsListView with multiple assignments
func (testSuite *AssignmentListViewTestSuite) TestAssignmentsListViewMultiple() {
	multipleAssignments := []clients.Assignment{
		{
			ID:   "test-id-1",
			Name: "Test Assignment 1",
			Due:  time.Now(),
			URL:  "https://example.com/1",
		},
		{
			ID:   "test-id-2",
			Name: "Test Assignment 2",
			Due:  time.Now().Add(24 * time.Hour),
			URL:  "https://example.com/2",
		},
	}
	embed := views.AssignmentsListView(testSuite.mockMember, multipleAssignments)

	testSuite.NotNil(embed)
	testSuite.Equal("assignments", embed.Title)
	testSuite.Equal("2 assignments in course", embed.Description)
	testSuite.Len(embed.Fields, 6) // 3 fields per assignment
}

// Test AssignmentsListActions
func (testSuite *AssignmentListViewTestSuite) TestAssignmentsListActions() {
	actions := views.AssignmentsListActions()

	testSuite.NotNil(actions)
	testSuite.Len(actions.Components, 1)

	button := actions.Components[0].(discordgo.Button)
	testSuite.Equal("➕", button.Emoji.Name)
	testSuite.Equal("add", button.Label)
	testSuite.Equal(discordgo.PrimaryButton, button.Style)
	testSuite.Equal("addAssignmentAction", button.CustomID)
}
