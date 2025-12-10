package views

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/stretchr/testify/suite"
)

// ConfigViewTestSuite tests config view functionality
type ConfigViewTestSuite struct {
	suite.Suite
}

func TestConfigView(t *testing.T) {
	suite.Run(t, new(ConfigViewTestSuite))
}

// Test ConfigView with empty course
func (s *ConfigViewTestSuite) TestConfigViewEmpty() {
	emptyCourse := clients.Course{}
	embed := ConfigView(emptyCourse)

	s.NotNil(embed)
	s.Equal("course config", embed.Title)
	s.Len(embed.Fields, 2)
	s.Equal("notifications channel", embed.Fields[0].Name)
	s.Equal("", embed.Fields[0].Value)
	s.Equal("notifications role", embed.Fields[1].Name)
	s.Equal("", embed.Fields[1].Value)
}

// Test ConfigView with channel only
func (s *ConfigViewTestSuite) TestConfigViewChannelOnly() {
	testCourse := clients.Course{
		CourseID:      "test-course",
		NotifyChannel: "123456789",
		NotifyGroup:   "",
	}
	embed := ConfigView(testCourse)

	s.NotNil(embed)
	s.Equal("course config", embed.Title)
	s.Len(embed.Fields, 2)
	s.Equal("notifications channel", embed.Fields[0].Name)
	s.Equal("<#123456789>", embed.Fields[0].Value)
	s.Equal("notifications role", embed.Fields[1].Name)
	s.Equal("", embed.Fields[1].Value)
}

// Test ConfigView with role only
func (s *ConfigViewTestSuite) TestConfigViewRoleOnly() {
	testCourse := clients.Course{
		CourseID:      "test-course",
		NotifyChannel: "",
		NotifyGroup:   "987654321",
	}
	embed := ConfigView(testCourse)

	s.NotNil(embed)
	s.Equal("course config", embed.Title)
	s.Len(embed.Fields, 2)
	s.Equal("notifications channel", embed.Fields[0].Name)
	s.Equal("", embed.Fields[0].Value)
	s.Equal("notifications role", embed.Fields[1].Name)
	s.Equal("<@&987654321>", embed.Fields[1].Value)
}

// Test ConfigView with both channel and role
func (s *ConfigViewTestSuite) TestConfigViewBoth() {
	testCourse := clients.Course{
		CourseID:      "test-course",
		NotifyChannel: "123456789",
		NotifyGroup:   "987654321",
	}
	embed := ConfigView(testCourse)

	s.NotNil(embed)
	s.Equal("course config", embed.Title)
	s.Len(embed.Fields, 2)
	s.Equal("notifications channel", embed.Fields[0].Name)
	s.Equal("<#123456789>", embed.Fields[0].Value)
	s.Equal("notifications role", embed.Fields[1].Name)
	s.Equal("<@&987654321>", embed.Fields[1].Value)
}

// Test ConfigActions
func (s *ConfigViewTestSuite) TestConfigActions() {
	components := ConfigActions()

	s.Len(components, 2)

	// Test channel select menu
	channelRow := components[0].(*discordgo.ActionsRow)
	s.Len(channelRow.Components, 1)
	channelMenu := channelRow.Components[0].(discordgo.SelectMenu)
	s.Equal(discordgo.ChannelSelectMenu, channelMenu.MenuType)
	s.Equal("updateNotifyChannel", channelMenu.CustomID)
	s.Equal("update notifications channel", channelMenu.Placeholder)

	// Test role select menu
	roleRow := components[1].(*discordgo.ActionsRow)
	s.Len(roleRow.Components, 1)
	roleMenu := roleRow.Components[0].(discordgo.SelectMenu)
	s.Equal(discordgo.RoleSelectMenu, roleMenu.MenuType)
	s.Equal("updateNotifyRole", roleMenu.CustomID)
	s.Equal("update notifications role", roleMenu.Placeholder)
}
