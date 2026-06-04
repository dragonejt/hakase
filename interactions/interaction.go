package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

type Interaction interface {
	UpdateAssignment(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
	UpdateAssignmentSubmit(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
	DeleteAssignment(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
	AddAssignment(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
	AddAssignmentSubmit(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
	UpdateNotifyChannel(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
	UpdateNotifyRole(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
	SlashAssignments(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
	SlashHakase(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
}

type InteractionHandler struct {
	Interaction
	HakaseClient clients.HakaseClient
}
