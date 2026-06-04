package clients

import "github.com/bwmarrin/discordgo"

type DiscordClient interface {
	UpdateCustomStatus(state string) (err error)
	Guild(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error)
	ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
	FollowupMessageCreate(interaction *discordgo.Interaction, wait bool, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (*discordgo.Message, error)
	UserGuilds(limit int, beforeID, afterID string, withCounts bool, options ...discordgo.RequestOption) (st []*discordgo.UserGuild, err error)
	GuildScheduledEvents(guildID string, userCount bool, options ...discordgo.RequestOption) (st []*discordgo.GuildScheduledEvent, err error)
	GuildScheduledEventEdit(guildID, eventID string, guildScheduledEventParams *discordgo.GuildScheduledEventParams, options ...discordgo.RequestOption) (st *discordgo.GuildScheduledEvent, err error)
	GuildScheduledEventCreate(guildID string, guildScheduledEventParams *discordgo.GuildScheduledEventParams, options ...discordgo.RequestOption) (st *discordgo.GuildScheduledEvent, err error)
	GuildScheduledEventDelete(guildID, eventID string, options ...discordgo.RequestOption) error
}
