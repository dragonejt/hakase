// Package events provides the Discord ready event handler.
package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

// Ready handles the Discord ready event and updates bot status and notifications.
func (handler *EventHandler) Ready(bot *discordgo.Session, ready *discordgo.Ready) {
	transaction := sentry.StartTransaction(context.Background(), "ready")
	defer transaction.Finish()
	slog.Info("logged in", "user", ready.User.String())

	err := bot.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to update status").Error())
	}
}
