// Package events provides the Discord interaction event handler.
package events

import (
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/palantir/stacktrace"
)

// InteractionCreate dispatches Discord interactions to the appropriate handler based on type and command.
func (handler *EventHandler) InteractionCreate(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	interaction := handler.InteractionHandler
	switch interactionCreate.Type {
	case discordgo.InteractionApplicationCommand:
		switch interactionCreate.ApplicationCommandData().Name {
		case "assignments":
			interaction.SlashAssignments(bot, interactionCreate)
		case "hakase":
			interaction.SlashHakase(bot, interactionCreate)
		default:
			slog.Error(stacktrace.NewError("unknown command: %s", interactionCreate.ApplicationCommandData().Name).Error())
		}
	case discordgo.InteractionMessageComponent:
		customID := interactionCreate.MessageComponentData().CustomID
		if strings.HasPrefix(customID, "addAssignmentAction") {
			interaction.AddAssignment(bot, interactionCreate)
		} else if strings.HasPrefix(customID, "updateAssignmentAction") {
			interaction.UpdateAssignment(bot, interactionCreate)
		} else if strings.HasPrefix(customID, "deleteAssignmentAction") {
			interaction.DeleteAssignment(bot, interactionCreate)
		} else if strings.HasPrefix(customID, "updateNotifyChannel") {
			interaction.UpdateNotifyChannel(bot, interactionCreate)
		} else if strings.HasPrefix(customID, "updateNotifyRole") {
			interaction.UpdateNotifyRole(bot, interactionCreate)
		} else {
			slog.Error(stacktrace.NewError("unknown message component action: %s", customID).Error())
		}
	case discordgo.InteractionModalSubmit:
		customID := interactionCreate.ModalSubmitData().CustomID
		if strings.HasPrefix(customID, "addAssignment") {
			interaction.AddAssignmentSubmit(bot, interactionCreate)
		} else if strings.HasPrefix(customID, "updateAssignment") {
			interaction.UpdateAssignmentSubmit(bot, interactionCreate)
		} else {
			slog.Error(stacktrace.NewError("unknown modal submit: %s", interactionCreate.ModalSubmitData().CustomID).Error())
		}
	default:
		slog.Error(stacktrace.NewError("unknown interaction type: %d", interactionCreate.Type).Error())

	}
}
