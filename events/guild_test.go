package events_test

import (
	"log/slog"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GuildEventsTestSuite struct {
	suite.Suite
	bot          *MockDiscordClient
	guildCreate  *discordgo.GuildCreate
	guildDelete  *discordgo.GuildDelete
	hakaseClient *MockHakaseClient
	event        *events.EventHandler
}

func TestGuildEvents(t *testing.T) {
	suite.Run(t, new(GuildEventsTestSuite))
}

type MockHakaseClient struct {
	clients.HakaseClient
	mock.Mock
}

func (hakaseClient *MockHakaseClient) CreateCourse(span *sentry.Span, course clients.Course) error {
	hakaseClient.Called(span, course)
	return nil
}

func (hakaseClient *MockHakaseClient) DeleteCourse(span *sentry.Span, courseID string) error {
	hakaseClient.Called(span, courseID)
	return nil
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

func (testSuite *GuildEventsTestSuite) SetupTest() {
	testSuite.bot = new(MockDiscordClient)
	guild := &discordgo.Guild{
		ID:   "1234567890",
		Name: "Test Guild",
	}
	userGuild := &discordgo.UserGuild{
		ID:   guild.ID,
		Name: guild.Name,
	}
	testSuite.guildCreate = &discordgo.GuildCreate{
		Guild: guild,
	}
	testSuite.guildDelete = &discordgo.GuildDelete{
		Guild:        guild,
		BeforeDelete: guild,
	}
	testSuite.hakaseClient = new(MockHakaseClient)
	testSuite.event = &events.EventHandler{
		HakaseClient: testSuite.hakaseClient,
	}
	testSuite.bot.On("UserGuilds", 100, "", "", false, []discordgo.RequestOption(nil)).Return([]*discordgo.UserGuild{userGuild}, nil)
	testSuite.bot.On("UpdateCustomStatus", mock.Anything).Return(nil)
}

func (testSuite *GuildEventsTestSuite) TestGuildCreateSuccess() {
	testSuite.hakaseClient.On("CreateCourse", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		slog.Info("CreateCourse called")
	})
	testSuite.event.GuildCreate(testSuite.bot, testSuite.guildCreate)
	testSuite.hakaseClient.AssertExpectations(testSuite.T())
	testSuite.bot.AssertExpectations(testSuite.T())
}

func (testSuite *GuildEventsTestSuite) TestGuildDeleteSuccess() {
	testSuite.hakaseClient.On("DeleteCourse", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		slog.Info("DeleteCourse called")
	})
	testSuite.event.GuildDelete(testSuite.bot, testSuite.guildDelete)
	testSuite.hakaseClient.AssertExpectations(testSuite.T())
	testSuite.bot.AssertExpectations(testSuite.T())
}
