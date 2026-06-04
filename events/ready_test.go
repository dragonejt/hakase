package events_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ReadyTestSuite tests ready event handling
type ReadyTestSuite struct {
	suite.Suite
	bot            *MockDiscordClient
	ready          *discordgo.Ready
	logBuffer      *bytes.Buffer
	originalLogger *slog.Logger
	event          *events.EventHandler
}

func TestReady(t *testing.T) {
	suite.Run(t, new(ReadyTestSuite))
}

func (testSuite *ReadyTestSuite) SetupTest() {
	testSuite.bot = new(MockDiscordClient)

	testSuite.ready = &discordgo.Ready{
		User: &discordgo.User{
			Username: "TestBot",
			ID:       "123456789",
		},
	}

	testSuite.logBuffer = new(bytes.Buffer)
	testSuite.originalLogger = slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(testSuite.logBuffer, nil)))
	testSuite.event = &events.EventHandler{}
	testSuite.bot.On("UserGuilds", 100, "", "", false, []discordgo.RequestOption(nil)).Return([]*discordgo.UserGuild{}, nil)
	testSuite.bot.On("UpdateCustomStatus", mock.Anything).Return(nil)
}

func (testSuite *ReadyTestSuite) TearDownTest() {
	slog.SetDefault(testSuite.originalLogger)
}

func (testSuite *ReadyTestSuite) TestReadyEventBasic() {
	testSuite.event.Ready(testSuite.bot, testSuite.ready)

	logOutput := testSuite.logBuffer.String()
	testSuite.Contains(logOutput, "user=TestBot")
	testSuite.bot.AssertExpectations(testSuite.T())
}
