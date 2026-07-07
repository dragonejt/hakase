package dev.dragonejt.hakase.clients

import dev.dragonejt.hakase.telemetry.LogBase
import dev.minn.jda.ktx.events.CoroutineEventListener
import dev.minn.jda.ktx.events.CoroutineEventManager
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.runBlocking
import net.dv8tion.jda.api.JDABuilder
import net.dv8tion.jda.api.hooks.EventListener
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@ConfigurationProperties(prefix = "discord") data class DiscordProperties(val token: String)

@Configuration
@EnableConfigurationProperties(DiscordProperties::class)
class DiscordBotConfiguration : LogBase() {
    @Bean
    fun discordBot(
        props: DiscordProperties,
        eventManager: CoroutineEventManager,
        eventListeners: Array<EventListener>,
        asyncEventListeners: Array<CoroutineEventListener>,
    ): JDABuilder = runBlocking {
        val bot = JDABuilder.createDefault(props.token)

        log.atDebug {
            message =
                "Registering Event Listeners: " +
                    arrayOf(*eventListeners, *asyncEventListeners).map { event ->
                        event.javaClass.simpleName
                    }
            payload =
                mapOf(
                    "event_listeners" to eventListeners,
                    "async_event_listeners" to asyncEventListeners,
                )
        }
        bot.addEventListeners(*eventListeners)
        bot.addEventListeners(*asyncEventListeners)
        bot.setEventManager(eventManager)

        return@runBlocking bot
    }

    @Bean
    fun scope(): CoroutineScope {
        return CoroutineScope(SupervisorJob() + Dispatchers.IO)
    }

    @Bean
    fun eventManager(scope: CoroutineScope): CoroutineEventManager {
        return CoroutineEventManager(scope)
    }
}
