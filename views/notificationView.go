package views

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

func NotificationView(assignment clients.Assignment, timeBuffer string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s is due in %s!", assignment.Name, timeBuffer),
		Description: fmt.Sprintf("due %s", assignment.Due.Format(time.RFC1123)),
		URL:         assignment.URL,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("id %s", assignment.ID)},
	}
}
