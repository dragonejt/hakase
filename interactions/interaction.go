package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/getsentry/sentry-go"
)

type Interaction interface {
	UpdateAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
	UpdateAssignmentSubmit(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
	DeleteAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
	AddAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
	AddAssignmentSubmit(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
	UpdateNotifyChannel(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
	UpdateNotifyRole(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
	SlashAssignments(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
	listAssignments(span *sentry.Span, interactionCreate *discordgo.InteractionCreate)
	SlashHakase(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

type InteractionHandler struct {
	Interaction
	HakaseClient clients.HakaseClient
	Bot          *discordgo.Session
}
