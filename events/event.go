package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/interactions"
)

type Event interface {
	GuildCreate(bot clients.DiscordClient, guildCreate *discordgo.GuildCreate)
	GuildDelete(bot clients.DiscordClient, guildDelete *discordgo.GuildDelete)
	InteractionCreate(bot clients.DiscordClient, interactionCreate *discordgo.InteractionCreate)
	Ready(bot clients.DiscordClient, ready *discordgo.Ready)
}

type EventHandler struct {
	HakaseClient       clients.HakaseClient
	InteractionHandler interactions.Interaction
}

func NewEventHandler(hakaseClient clients.HakaseClient) *EventHandler {
	return &EventHandler{
		HakaseClient:       hakaseClient,
		InteractionHandler: &interactions.InteractionHandler{HakaseClient: hakaseClient},
	}
}
