// Package interactions provides handlers for the /hakase slash command.
package interactions

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/views"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

var HakaseCommand = discordgo.ApplicationCommand{
	Name:        "hakase",
	Description: "course configuration",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "cmd",
			Description: "subcommand to execute",
			Type:        discordgo.ApplicationCommandOptionString,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "rock-paper-scissors",
					Value: "rock-paper-scissors",
				},
				{
					Name:  "config",
					Value: "config",
				},
			},
		},
	},
}

var rockPaperScissorsGIFS = []string{
	"https://tenor.com/view/hakase-rock-paper-scissors-nichijou-gif-16268712187530688616",
	"https://tenor.com/view/nichijou-hakase-rps-rock-paper-scissors-nano-gif-17854309283562565671",
	"https://tenor.com/view/rps-nichijou-hakase-nano-rock-paper-scissors-gif-9851171842395079248",
	"https://tenor.com/view/hakase-rock-paper-scissors-nichijou-gif-8852988661228140380",
	"https://tenor.com/view/rps-nichijou-hakase-nano-rock-paper-scissors-gif-9851171842395079248",
	"https://tenor.com/view/nichijou-hakase-rps-rock-paper-scissors-nano-gif-11850067363499322337",
}

// SlashHakase handles the /hakase slash command interaction.
// It dispatches subcommands such as rock-paper-scissors and config.
func (handler *InteractionHandler) SlashHakase(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(interactionCreate.ApplicationCommandData().Options))
	for _, opt := range interactionCreate.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	slog.Info(fmt.Sprintf("/hakase executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	transaction := sentry.StartTransaction(context.Background(), "/hakase")
	defer transaction.Finish()

	subcommand, exists := optionMap["cmd"]
	if !exists {
		handler.ping(transaction, interactionCreate)
	} else {
		switch subcommand.StringValue() {
		case "rock-paper-scissors":
			handler.rockPaperScissors(transaction, interactionCreate)
		case "config":
			handler.config(transaction, interactionCreate)
		}
	}
}

// ping responds to the /hakase command with a pong and backend response time.
func (handler *InteractionHandler) ping(span *sentry.Span, interactionCreate *discordgo.InteractionCreate) {
	span = span.StartChild("/hakase ping")
	defer span.Finish()

	start := time.Now()
	_, err := handler.HakaseClient.ReadCourse(span, interactionCreate.GuildID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error pinging backend").Error())
	}

	err = handler.Bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("hakase pong! response time: %dms", time.Since(start).Milliseconds()),
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}

}

// rockPaperScissors responds with a random rock-paper-scissors GIF.
func (handler *InteractionHandler) rockPaperScissors(span *sentry.Span, interactionCreate *discordgo.InteractionCreate) {
	span = span.StartChild("/hakase rockPaperScissors")
	defer span.Finish()

	err := handler.Bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: rockPaperScissorsGIFS[rand.Intn(len(rockPaperScissorsGIFS))],
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}
}

// config responds with the course configuration embed and components.
func (handler *InteractionHandler) config(span *sentry.Span, interactionCreate *discordgo.InteractionCreate) {
	span = span.StartChild("/hakase config")
	defer span.Finish()

	course, err := handler.HakaseClient.ReadCourse(span, interactionCreate.GuildID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error reading course").Error())
		err = handler.Bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("error reading course: %s", err.Error()),
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	err = handler.Bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{views.ConfigView(course)},
			Components: views.ConfigActions(),
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}
}
