// Package interactions provides handlers for the /assignments slash command.
package interactions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/views"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

var AssignmentsCommand = discordgo.ApplicationCommand{
	Name:        "assignments",
	Description: "configure assignments for due date notifications",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "id",
			Description: "retrieves assignment with this id",
			Type:        discordgo.ApplicationCommandOptionString,
		},
	},
}

// SlashAssignments handles the /assignments slash command interaction.
// It retrieves a specific assignment or lists all assignments for the guild.
func (handler *InteractionHandler) SlashAssignments(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(interactionCreate.ApplicationCommandData().Options))
	for _, opt := range interactionCreate.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	slog.Info(fmt.Sprintf("/assignments executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	transaction := sentry.StartTransaction(context.Background(), "/assignments")
	defer transaction.Finish()

	assignmentID, exists := optionMap["id"]
	if exists {
		handler.getAssignment(transaction, interactionCreate, assignmentID.StringValue())
	} else {
		handler.listAssignments(transaction, interactionCreate)
	}

}

// getAssignment retrieves and responds with a specific assignment's details.
func (handler *InteractionHandler) getAssignment(span *sentry.Span, interactionCreate *discordgo.InteractionCreate, assignmentID string) {
	span = span.StartChild("/assignments getAssignment")
	defer span.Finish()

	assignment, err := handler.HakaseClient.ReadAssignment(span, assignmentID)

	if err != nil {
		err = handler.Bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	} else {
		err = handler.Bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{views.AssignmentView(interactionCreate.Member, assignment)},
				Components: []discordgo.MessageComponent{views.AssignmentActions(assignment)},
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	}
}

// listAssignments retrieves and responds with a list of assignments for the guild.
func (handler *InteractionHandler) listAssignments(span *sentry.Span, interactionCreate *discordgo.InteractionCreate) {
	span = span.StartChild("/assignments listAssignments")
	defer span.Finish()

	assignments, err := handler.HakaseClient.ListAssignments(span, interactionCreate.GuildID)
	if err != nil {
		err = handler.Bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	} else {
		err = handler.Bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{views.AssignmentsListView(interactionCreate.Member, assignments)},
				Components: []discordgo.MessageComponent{views.AssignmentsListActions()},
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	}
}
