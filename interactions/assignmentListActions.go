// Package interactions provides handlers for assignment list actions (add assignment).
package interactions

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

// AddAssignment opens a modal for adding a new assignment via Discord interaction.
func (handler *InteractionHandler) AddAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	transaction := sentry.StartTransaction(context.Background(), "addAssignmentAction")
	defer transaction.Finish()
	slog.Debug("addAssignment executed", "username", interactionCreate.Member.User.Username, "user_id", interactionCreate.Member.User.ID, "guild_id", interactionCreate.GuildID)
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

	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   "addAssignment",
			Title:      "add assignment",
			Components: views.AssignmentModal(nil),
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}

}

// AddAssignmentSubmit handles the submission of the add assignment modal and creates the assignment.
func (handler *InteractionHandler) AddAssignmentSubmit(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	slog.Info("addAssignmentSubmit executed", "username", interactionCreate.Member.User.Username, "user_id", interactionCreate.Member.User.ID, "guild_id", interactionCreate.GuildID)
	transaction := sentry.StartTransaction(context.Background(), "addAssignmentSubmit")
	defer transaction.Finish()

	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}

	assignmentData := interactionCreate.ModalSubmitData()
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

	assignment := clients.Assignment{
		Name:     assignmentData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		Due:      due,
		CourseID: interactionCreate.GuildID,
	}

	if assignmentData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value != "" {
		assignment.URL = assignmentData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	if assignment.Due.Before(time.Now()) {
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: "due date before current time! hakase does not support this.",
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	createdAssignmentID, err := handler.HakaseClient.CreateAssignment(transaction, assignment)
	if err != nil {
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: err.Error(),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}
	createdAssignment, err := handler.HakaseClient.ReadAssignment(transaction, createdAssignmentID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error getting created assignment with id: %s", createdAssignmentID).Error())
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: err.Error(),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	startTime := assignment.Due.Add(-1 * time.Hour)
	//nolint:errcheck
	go bot.GuildScheduledEventCreate(assignment.CourseID, &discordgo.GuildScheduledEventParams{
		Name:               assignment.Name,
		Description:        assignment.ID,
		ScheduledStartTime: &startTime,
		ScheduledEndTime:   &assignment.Due,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		EntityMetadata: &discordgo.GuildScheduledEventEntityMetadata{
			Location: assignment.URL,
		},
	})

	_, err = bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
		Content:    "assignment created!",
		Embeds:     []*discordgo.MessageEmbed{views.AssignmentView(interactionCreate.Member, createdAssignment)},
		Components: []discordgo.MessageComponent{views.AssignmentActions(createdAssignment)},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}
}
