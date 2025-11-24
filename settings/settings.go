// Package settings provides environment variable configuration for the bot.
package settings

import "os"

var ENV string = os.Getenv("ENV")
var DEBUG bool = ENV != "production"
var DISCORD_BOT_TOKEN string = os.Getenv("DISCORD_BOT_TOKEN")
var SENTRY_DSN string = os.Getenv("SENTRY_DSN")
var CF_ACCOUNT_ID = os.Getenv("CF_ACCOUNT_ID")
var CF_API_TOKEN = os.Getenv("CF_API_TOKEN")
var D1_DATABASE_ID = os.Getenv("D1_DATABASE_ID")
