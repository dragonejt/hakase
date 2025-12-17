// Package settings provides environment variable configuration for the bot.
package settings

import "os"

var ENV string = os.Getenv("ENV")
var DEBUG bool = ENV != "production"
var DISCORD_BOT_TOKEN string = os.Getenv("DISCORD_BOT_TOKEN")
var SENTRY_DSN string = os.Getenv("SENTRY_DSN")
var BACKEND_URL string = os.Getenv("BACKEND_URL")
var BACKEND_AUTH_TOKEN string = os.Getenv("BACKEND_AUTH_TOKEN")
