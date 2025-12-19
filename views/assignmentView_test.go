package views_test

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
	"github.com/stretchr/testify/suite"
)

// AssignmentViewTestSuite tests assignment view functionality
type AssignmentViewTestSuite struct {
	suite.Suite
	mockMember *discordgo.Member
}

func TestAssignmentView(testSuite *testing.T) {
	suite.Run(testSuite, new(AssignmentViewTestSuite))
}

func (s *AssignmentViewTestSuite) SetupTest() {
	s.mockMember = &discordgo.Member{
		User: &discordgo.User{
			Username: "TestUser",
			ID:       "123456789",
		},
	}
}

// Test AssignmentView
func (s *AssignmentViewTestSuite) TestAssignmentView() {
	testTime := time.Now()
	testAssignment := clients.Assignment{
		ID:       "test-id-1",
		Name:     "Test Assignment",
		Due:      &testTime,
		URL:      "https://example.com/assignment",
		CourseID: "test-course",
	}

	embed := views.AssignmentView(s.mockMember, testAssignment)

	s.NotNil(embed)
	s.Equal("Test Assignment", embed.Title)
	s.Equal("due "+testTime.Format(time.RFC1123), embed.Description)
	s.Equal("TestUser", embed.Author.Name)
	s.Equal("https://example.com/assignment", embed.URL)
	s.Equal("id test-id-1", embed.Footer.Text)
}

// Test AssignmentActions
func (s *AssignmentViewTestSuite) TestAssignmentActions() {
	testAssignment := clients.Assignment{
		ID:   "test-id-1",
		Name: "Test Assignment",
	}

	actions := views.AssignmentActions(testAssignment)

	s.NotNil(actions)
	s.Len(actions.Components, 2)

	// Test edit button
	editButton := actions.Components[0].(discordgo.Button)
	s.Equal("📝", editButton.Emoji.Name)
	s.Equal("edit", editButton.Label)
	s.Equal(discordgo.PrimaryButton, editButton.Style)
	s.Equal("updateAssignmentAction_test-id-1", editButton.CustomID)

	// Test remove button
	removeButton := actions.Components[1].(discordgo.Button)
	s.Equal("🗑️", removeButton.Emoji.Name)
	s.Equal("remove", removeButton.Label)
	s.Equal(discordgo.SecondaryButton, removeButton.Style)
	s.Equal("deleteAssignmentAction_test-id-1", removeButton.CustomID)
}

// Test AssignmentModal for new assignment
func (s *AssignmentViewTestSuite) TestAssignmentModalNew() {
	components := views.AssignmentModal(nil)

	s.Len(components, 3) // Should have 3 action rows

	// Test name input
	nameRow := components[0].(discordgo.ActionsRow)
	s.Len(nameRow.Components, 1)
	nameInput := nameRow.Components[0].(discordgo.TextInput)
	s.Equal("assignmentName", nameInput.CustomID)
	s.Equal("assignment name:", nameInput.Label)
	s.Equal(discordgo.TextInputShort, nameInput.Style)
	s.Equal("Assignment 1", nameInput.Placeholder)
	s.True(nameInput.Required)
	s.Equal(50, nameInput.MaxLength)

	// Test due date input
	dueRow := components[1].(discordgo.ActionsRow)
	s.Len(dueRow.Components, 1)
	dueInput := dueRow.Components[0].(discordgo.TextInput)
	s.Equal("assignmentDue", dueInput.CustomID)
	s.Equal("due date:", dueInput.Label)
	s.Equal(discordgo.TextInputShort, dueInput.Style)
	s.True(dueInput.Required)

	// Test link input
	linkRow := components[2].(discordgo.ActionsRow)
	s.Len(linkRow.Components, 1)
	linkInput := linkRow.Components[0].(discordgo.TextInput)
	s.Equal("assignmentLink", linkInput.CustomID)
	s.Equal("link:", linkInput.Label)
	s.Equal(discordgo.TextInputShort, linkInput.Style)
	s.Equal("https://canvas.instructure.com", linkInput.Placeholder)
	s.False(linkInput.Required)
}

// Test AssignmentModal for existing assignment
func (s *AssignmentViewTestSuite) TestAssignmentModalExisting() {
	testTime := time.Now()
	testAssignment := &clients.Assignment{
		ID:   "test-id-1",
		Name: "Existing Assignment",
		Due:  &testTime,
		URL:  "https://example.com/existing",
	}

	components := views.AssignmentModal(testAssignment)

	s.Len(components, 3)

	// Test that existing assignment values are used as placeholders
	nameRow := components[0].(discordgo.ActionsRow)
	nameInput := nameRow.Components[0].(discordgo.TextInput)
	s.Equal("Existing Assignment", nameInput.Placeholder)
	s.False(nameInput.Required) // Should not be required for existing assignment

	dueRow := components[1].(discordgo.ActionsRow)
	dueInput := dueRow.Components[0].(discordgo.TextInput)
	s.Equal(testTime.Format(time.RFC1123), dueInput.Placeholder)
	s.False(dueInput.Required)

	linkRow := components[2].(discordgo.ActionsRow)
	linkInput := linkRow.Components[0].(discordgo.TextInput)
	s.Equal("https://example.com/existing", linkInput.Placeholder)
	s.False(linkInput.Required)
}
