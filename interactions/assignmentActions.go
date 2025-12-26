// Package interactions provides handlers for assignment actions (edit, update, delete).
package interactions

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/views"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

// UpdateAssignment opens a modal for updating an assignment via Discord interaction.
func (handler *InteractionHandler) UpdateAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	transaction := sentry.StartTransaction(context.Background(), "updateAssignmentAction")
	defer transaction.Finish()
	slog.Debug("updateAssignment executed", "username", interactionCreate.Member.User.Username, "user_id", interactionCreate.Member.User.ID, "guild_id", interactionCreate.GuildID)
	if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "admin permissions needed!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	assignmentID := strings.Split(interactionCreate.MessageComponentData().CustomID, "_")[1]
	assignment, err := handler.HakaseClient.ReadAssignment(transaction, assignmentID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error reading assignment").Error())
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}
	err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   fmt.Sprintf("updateAssignment_%s", assignmentID),
			Title:      "update assignment",
			Components: views.AssignmentModal(&assignment),
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}
}

// UpdateAssignmentSubmit handles the submission of the update assignment modal and updates the assignment.
func (handler *InteractionHandler) UpdateAssignmentSubmit(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	slog.Info("updateAssignmentSubmit executed", "username", interactionCreate.Member.User.Username, "user_id", interactionCreate.Member.User.ID, "guild_id", interactionCreate.GuildID)
	transaction := sentry.StartTransaction(context.Background(), "updateAssignmentSubmit")
	defer transaction.Finish()

	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}

	assignmentID := strings.Split(interactionCreate.ModalSubmitData().CustomID, "_")[1]
	assignmentData := interactionCreate.ModalSubmitData()
	assignment, err := handler.HakaseClient.ReadAssignment(transaction, assignmentID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error getting current assignment with id: %s", assignmentID).Error())
	}

	if assignmentData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value != "" {
		due, err := dateparse.ParseAny(assignmentData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
		if err != nil {
			_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
				Content: fmt.Sprintf("error parsing due date: %s", err.Error()),
			})
			if err != nil {
				slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
			}
			return
		}
		assignment.Due = due
	}

	if assignmentData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value != "" {
		assignment.Name = assignmentData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	if assignmentData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value != "" {
		assignment.URL = assignmentData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	currentAssignment, err := handler.HakaseClient.ReadAssignment(transaction, assignmentID)
	if assignment.Due.Equal(time.Time{}) {
		assignment.Due = currentAssignment.Due
	} else if err == nil && assignment.Due.Before(currentAssignment.Due) {
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: "new due date before original assignment due date! hakase does not support this.",
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	assignment.Assignment = currentAssignment.ID
	err = handler.HakaseClient.UpdateAssignment(transaction, assignment)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error updating assignment").Error())
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: err.Error(),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}
	updatedAssignment, err := handler.HakaseClient.ReadAssignment(transaction, assignmentID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error reading updated assignment").Error())
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: err.Error(),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	_, err = bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
		Content:    "assignment updated!",
		Embeds:     []*discordgo.MessageEmbed{views.AssignmentView(interactionCreate.Member, updatedAssignment)},
		Components: []discordgo.MessageComponent{views.AssignmentActions(updatedAssignment)},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}
}

// DeleteAssignment deletes an assignment based on user interaction.
func (handler *InteractionHandler) DeleteAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	slog.Debug("deleteAssignment executed", "username", interactionCreate.Member.User.Username, "user_id", interactionCreate.Member.User.ID, "guild_id", interactionCreate.GuildID)
	transaction := sentry.StartTransaction(context.Background(), "deleteAssignmentAction")
	defer transaction.Finish()

	assignmentID := strings.Split(interactionCreate.MessageComponentData().CustomID, "_")[1]
	err := handler.HakaseClient.DeleteAssignment(transaction, assignmentID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "unable to delete assignment %s", assignmentID).Error())
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("unable to delete assignment %s: %s", assignmentID, err.Error()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("assignment %s deleted!", assignmentID),
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}
}
