// hakase-discord is the entry point for the Discord bot.
// It initializes logging, Sentry, Discord session, backend client, and event handlers.
// It registers application commands and starts the bot event loop.
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/dragonejt/hakase-discord/interactions"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

func main() {
	if os.Getenv("ENV") != "production" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if sentryDsn := os.Getenv("SENTRY_DSN"); sentryDsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDsn,
			EnableTracing:    true,
			SampleRate:       1,
			TracesSampleRate: 1,
			SendDefaultPII:   true,
			EnableLogs:       true,
			Environment:      os.Getenv("ENV"),
		})
		if err != nil {
			slog.Warn(stacktrace.Propagate(err, "failed to initiate sentry").Error())
		}
		slog.SetDefault(slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, sentry.NewLogger(context.Background())), &slog.HandlerOptions{AddSource: true})))
	}

	bot, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_BOT_TOKEN")))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to create discord session").Error())
		return
	}
	bot.StateEnabled = true

	hakaseClient := &clients.BackendClient{
		URL:        os.Getenv("BACKEND_URL"),
		AuthToken:  os.Getenv("BACKEND_AUTH_TOKEN"),
		HTTPClient: bot.Client,
	}
	event := events.NewEventHandler(hakaseClient)

	err = bot.Open()
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to open discord session").Error())
		return
	}
	defer bot.Close() //nolint:errcheck

	bot.StateEnabled = true

	slog.Info("registering event handlers")
	bot.AddHandler(event.Ready)
	bot.AddHandler(event.GuildCreate)
	bot.AddHandler(event.GuildDelete)
	bot.AddHandler(event.InteractionCreate)

	slog.Info("registering slash commands")
	commands := []*discordgo.ApplicationCommand{&interactions.AssignmentsCommand, &interactions.HakaseCommand}
	for _, cmd := range commands {
		_, err = bot.ApplicationCommandCreate(bot.State.User.ID, "", cmd)
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "failed to register command: %s", cmd.Name).Error())
		} else {
			slog.Info("successfully registered command", "command_name", cmd.Name)
		}
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go event.RegisterAssignmentHandler(bot, shutdown)

	<-shutdown
}
