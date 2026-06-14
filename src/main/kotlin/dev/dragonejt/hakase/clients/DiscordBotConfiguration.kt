package dev.dragonejt.hakase.clients

import dev.dragonejt.hakase.events.EventHandler
import dev.kord.core.Kord
import dev.kord.core.event.gateway.GatewayEvent
import kotlinx.coroutines.runBlocking
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@ConfigurationProperties(prefix = "discord") data class DiscordProperties(val token: String)

@Configuration
@EnableConfigurationProperties(DiscordProperties::class)
class DiscordBotConfiguration(
    private val kordFactory: (String) -> Kord = { token -> runBlocking { Kord(token) } }
) {
    @Bean
    fun discordBot(
        properties: DiscordProperties,
        eventHandlers: List<EventHandler<out GatewayEvent>>,
    ): Kord {
        val bot: Kord = kordFactory(properties.token)

        eventHandlers.forEach { handler -> handler.register(bot) }
        return bot
    }
}
