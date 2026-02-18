package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/interactions"
)

type Event interface {
	GuildCreate(bot *discordgo.Session, guildCreate *discordgo.GuildCreate)
	GuildDelete(bot *discordgo.Session, guildDelete *discordgo.GuildDelete)
	InteractionCreate(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
	Ready(bot *discordgo.Session, ready *discordgo.Ready)
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
