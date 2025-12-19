package events

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/interactions"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type InteractionsTestSuite struct {
	suite.Suite
	mockInteraction *MockInteraction
	eventHandler    *EventHandler
}

func TestInteractions(t *testing.T) {
	suite.Run(t, new(InteractionsTestSuite))
}

type MockInteraction struct {
	interactions.Interaction
	mock.Mock
}

func (m *MockInteraction) UpdateAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) UpdateAssignmentSubmit(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) DeleteAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) AddAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) AddAssignmentSubmit(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) UpdateNotifyChannel(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) UpdateNotifyRole(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) SlashAssignments(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) SlashHakase(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (suite *InteractionsTestSuite) SetupTest() {
	suite.mockInteraction = &MockInteraction{}
	suite.eventHandler = &EventHandler{
		InteractionHandler: suite.mockInteraction,
	}
}

func (suite *InteractionsTestSuite) TestInteractionCreate_ApplicationCommand_Assignments() {
	bot := &discordgo.Session{}
	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "assignments",
			},
		},
	}

	suite.mockInteraction.On("SlashAssignments", bot, interactionCreate).Return()

	suite.eventHandler.InteractionCreate(bot, interactionCreate)

	suite.mockInteraction.AssertExpectations(suite.T())
}

func (suite *InteractionsTestSuite) TestInteractionCreate_ApplicationCommand_Hakase() {
	bot := &discordgo.Session{}
	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "hakase",
			},
		},
	}

	suite.mockInteraction.On("SlashHakase", bot, interactionCreate).Return()

	suite.eventHandler.InteractionCreate(bot, interactionCreate)

	suite.mockInteraction.AssertExpectations(suite.T())
}

func (suite *InteractionsTestSuite) TestInteractionCreate_MessageComponent_AddAssignment() {
	bot := &discordgo.Session{}
	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
			Data: discordgo.MessageComponentInteractionData{
				CustomID: "addAssignmentAction-test",
			},
		},
	}

	suite.mockInteraction.On("AddAssignment", bot, interactionCreate).Return()

	suite.eventHandler.InteractionCreate(bot, interactionCreate)

	suite.mockInteraction.AssertExpectations(suite.T())
}

func (suite *InteractionsTestSuite) TestInteractionCreate_MessageComponent_UpdateAssignment() {
	bot := &discordgo.Session{}
	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
			Data: discordgo.MessageComponentInteractionData{
				CustomID: "updateAssignmentAction-test",
			},
		},
	}

	suite.mockInteraction.On("UpdateAssignment", bot, interactionCreate).Return()

	suite.eventHandler.InteractionCreate(bot, interactionCreate)

	suite.mockInteraction.AssertExpectations(suite.T())
}

func (suite *InteractionsTestSuite) TestInteractionCreate_MessageComponent_DeleteAssignment() {
	bot := &discordgo.Session{}
	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
			Data: discordgo.MessageComponentInteractionData{
				CustomID: "deleteAssignmentAction-test",
			},
		},
	}

	suite.mockInteraction.On("DeleteAssignment", bot, interactionCreate).Return()

	suite.eventHandler.InteractionCreate(bot, interactionCreate)

	suite.mockInteraction.AssertExpectations(suite.T())
}

func (suite *InteractionsTestSuite) TestInteractionCreate_MessageComponent_UpdateNotifyChannel() {
	bot := &discordgo.Session{}
	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
			Data: discordgo.MessageComponentInteractionData{
				CustomID: "updateNotifyChannel-test",
			},
		},
	}

	suite.mockInteraction.On("UpdateNotifyChannel", bot, interactionCreate).Return()

	suite.eventHandler.InteractionCreate(bot, interactionCreate)

	suite.mockInteraction.AssertExpectations(suite.T())
}

func (suite *InteractionsTestSuite) TestInteractionCreate_MessageComponent_UpdateNotifyRole() {
	bot := &discordgo.Session{}
	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
			Data: discordgo.MessageComponentInteractionData{
				CustomID: "updateNotifyRole-test",
			},
		},
	}

	suite.mockInteraction.On("UpdateNotifyRole", bot, interactionCreate).Return()

	suite.eventHandler.InteractionCreate(bot, interactionCreate)

	suite.mockInteraction.AssertExpectations(suite.T())
}

func (suite *InteractionsTestSuite) TestInteractionCreate_ModalSubmit_AddAssignment() {
	bot := &discordgo.Session{}
	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionModalSubmit,
			Data: discordgo.ModalSubmitInteractionData{
				CustomID: "addAssignment-test",
			},
		},
	}

	suite.mockInteraction.On("AddAssignmentSubmit", bot, interactionCreate).Return()

	suite.eventHandler.InteractionCreate(bot, interactionCreate)

	suite.mockInteraction.AssertExpectations(suite.T())
}

func (suite *InteractionsTestSuite) TestInteractionCreate_ModalSubmit_UpdateAssignment() {
	bot := &discordgo.Session{}
	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionModalSubmit,
			Data: discordgo.ModalSubmitInteractionData{
				CustomID: "updateAssignment-test",
			},
		},
	}

	suite.mockInteraction.On("UpdateAssignmentSubmit", bot, interactionCreate).Return()

	suite.eventHandler.InteractionCreate(bot, interactionCreate)

	suite.mockInteraction.AssertExpectations(suite.T())
}
