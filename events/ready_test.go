package events

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/stretchr/testify/suite"
)

// ReadyTestSuite tests the ready event handling
type ReadyTestSuite struct {
	suite.Suite
	mockBot          *discordgo.Session
	mockReady        *discordgo.Ready
	mockHakaseClient clients.HakaseClient
	logBuffer        *bytes.Buffer
	originalLogger   *slog.Logger
}

func TestReady(t *testing.T) {
	suite.Run(t, new(ReadyTestSuite))
}

func (s *ReadyTestSuite) SetupTest() {
	s.mockBot = &discordgo.Session{}
	s.mockBot.State = &discordgo.State{}

	// Create a mock ready event
	s.mockReady = &discordgo.Ready{
		User: &discordgo.User{
			Username: "TestBot",
			ID:       "123456789",
		},
	}

	s.mockHakaseClient = nil // Not needed for ready event

	// Set up log capture
	s.logBuffer = new(bytes.Buffer)
	s.originalLogger = slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(s.logBuffer, nil)))
}

func (s *ReadyTestSuite) TearDownTest() {
	// Restore original logger
	slog.SetDefault(s.originalLogger)
}

// Test that Ready function handles the event without panicking
func (s *ReadyTestSuite) TestReadyEventBasic() {
	// This test verifies that the Ready function can handle a basic ready event
	// The function should update the bot status and log the login
	Ready(s.mockBot, s.mockReady, s.mockHakaseClient)

	// Verify that the login was logged
	logOutput := s.logBuffer.String()
	s.Contains(logOutput, "logged in as TestBot")

	// No error expected - function should complete without panic
}

// Test Ready function with empty guilds list
func (s *ReadyTestSuite) TestReadyEventNoGuilds() {
	// Test with no guilds - should handle gracefully
	s.mockBot.State.Guilds = []*discordgo.Guild{}
	Ready(s.mockBot, s.mockReady, s.mockHakaseClient)

	// Verify that login was logged and status update was attempted
	logOutput := s.logBuffer.String()
	s.Contains(logOutput, "logged in as TestBot")
	s.Contains(logOutput, "failed to update status") // Status update fails in test environment
}

// Test Ready function with multiple guilds
func (s *ReadyTestSuite) TestReadyEventMultipleGuilds() {
	// Test with multiple guilds
	s.mockBot.State.Guilds = []*discordgo.Guild{
		{ID: "guild1", Name: "Test Guild 1"},
		{ID: "guild2", Name: "Test Guild 2"},
		{ID: "guild3", Name: "Test Guild 3"},
	}
	Ready(s.mockBot, s.mockReady, s.mockHakaseClient)

	// Verify that login was logged and status update was attempted
	logOutput := s.logBuffer.String()
	s.Contains(logOutput, "logged in as TestBot")
	s.Contains(logOutput, "failed to update status") // Status update fails in test environment
}
