package events

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
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

func (m *MockInteraction) UpdateAssignment(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) UpdateAssignmentSubmit(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) DeleteAssignment(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) AddAssignment(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) AddAssignmentSubmit(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) UpdateNotifyChannel(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) UpdateNotifyRole(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) SlashAssignments(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

func (m *MockInteraction) SlashHakase(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate) {
	m.Called(bot, interactionCreate)
}

type MockDiscordClient struct {
	mock.Mock
}

func (m *MockDiscordClient) UpdateCustomStatus(state string) error {
	args := m.Called(state)
	return args.Error(0)
}

func (m *MockDiscordClient) Guild(guildID string, options ...discordgo.RequestOption) (*discordgo.Guild, error) {
	args := m.Called(guildID, options)
	return args.Get(0).(*discordgo.Guild), args.Error(1)
}

func (m *MockDiscordClient) ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	args := m.Called(channelID, data, options)
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

func (m *MockDiscordClient) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error {
	args := m.Called(interaction, resp, options)
	return args.Error(0)
}

func (m *MockDiscordClient) FollowupMessageCreate(interaction *discordgo.Interaction, wait bool, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	args := m.Called(interaction, wait, data, options)
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

func (m *MockDiscordClient) UserGuilds(limit int, beforeID, afterID string, withCounts bool, options ...discordgo.RequestOption) ([]*discordgo.UserGuild, error) {
	args := m.Called(limit, beforeID, afterID, withCounts, options)
	return args.Get(0).([]*discordgo.UserGuild), args.Error(1)
}

func (m *MockDiscordClient) GuildScheduledEvents(guildID string, userCount bool, options ...discordgo.RequestOption) ([]*discordgo.GuildScheduledEvent, error) {
	args := m.Called(guildID, userCount, options)
	return args.Get(0).([]*discordgo.GuildScheduledEvent), args.Error(1)
}

func (m *MockDiscordClient) GuildScheduledEventEdit(guildID, eventID string, params *discordgo.GuildScheduledEventParams, options ...discordgo.RequestOption) (*discordgo.GuildScheduledEvent, error) {
	args := m.Called(guildID, eventID, params, options)
	return args.Get(0).(*discordgo.GuildScheduledEvent), args.Error(1)
}

func (m *MockDiscordClient) GuildScheduledEventCreate(guildID string, params *discordgo.GuildScheduledEventParams, options ...discordgo.RequestOption) (*discordgo.GuildScheduledEvent, error) {
	args := m.Called(guildID, params, options)
	return args.Get(0).(*discordgo.GuildScheduledEvent), args.Error(1)
}

func (m *MockDiscordClient) GuildScheduledEventDelete(guildID, eventID string, options ...discordgo.RequestOption) error {
	args := m.Called(guildID, eventID, options)
	return args.Error(0)
}

func (suite *InteractionsTestSuite) SetupTest() {
	suite.mockInteraction = &MockInteraction{}
	suite.eventHandler = &EventHandler{
		InteractionHandler: suite.mockInteraction,
	}
}

func (suite *InteractionsTestSuite) TestInteractionCreate_ApplicationCommand_Assignments() {
	bot := new(MockDiscordClient)
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
	bot := new(MockDiscordClient)
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
	bot := new(MockDiscordClient)
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
	bot := new(MockDiscordClient)
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
	bot := new(MockDiscordClient)
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
	bot := new(MockDiscordClient)
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
	bot := new(MockDiscordClient)
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
	bot := new(MockDiscordClient)
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
	bot := new(MockDiscordClient)
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
