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
class DiscordBotConfiguration {

  @Bean
  fun kord(
      properties: DiscordProperties,
      eventHandlers: List<EventHandler<out GatewayEvent>>,
  ): Kord = runBlocking { discordBot(properties, eventHandlers) }

  suspend fun discordBot(
      properties: DiscordProperties,
      eventHandlers: List<EventHandler<out GatewayEvent>>,
  ): Kord {
    val bot = Kord(properties.token)

    eventHandlers.forEach { handler -> handler.register(bot) }
    return bot
  }
}
