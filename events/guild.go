// Package events provides Discord event handlers for guild create and delete.
package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

// GuildCreate handles the event when the bot is added to a guild and creates a course.
func (handler *EventHandler) GuildCreate(bot *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	transaction := sentry.StartTransaction(context.Background(), "guildCreate")
	defer transaction.Finish()
	slog.Info("added to guild", "guild_name", guildCreate.Name, "guild_id", guildCreate.ID)

	course := clients.Course{
		CourseID: guildCreate.ID,
	}
	err := handler.HakaseClient.CreateCourse(transaction, course)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to create course").Error())
	}

	err = bot.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to update status").Error())
	}
}

// GuildDelete handles the event when the bot is removed from a guild and deletes the course.
func (handler *EventHandler) GuildDelete(bot *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	transaction := sentry.StartTransaction(context.Background(), "guildDelete")
	defer transaction.Finish()
	slog.Info("removed from guild", "guild_name", guildDelete.Name, "guild_id", guildDelete.ID)

	err := handler.HakaseClient.DeleteCourse(transaction, guildDelete.ID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to delete course").Error())
	}

	err = bot.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to update status").Error())
	}
}
