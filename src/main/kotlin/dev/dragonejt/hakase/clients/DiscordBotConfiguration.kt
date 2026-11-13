package dev.dragonejt.hakase.clients

import dev.minn.jda.ktx.events.CoroutineEventListener
import dev.minn.jda.ktx.events.CoroutineEventManager
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.runBlocking
import net.dv8tion.jda.api.JDABuilder
import net.dv8tion.jda.api.hooks.EventListener
import org.apache.logging.log4j.kotlin.Logging
import org.apache.logging.log4j.message.StringMapMessage
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@ConfigurationProperties(prefix = "discord")
data class DiscordProperties(val token: String, val rpsGifs: List<String>)

@Configuration
@EnableConfigurationProperties(DiscordProperties::class)
class DiscordBotConfiguration {
    companion object : Logging

    @Bean
    fun discordBot(
        props: DiscordProperties,
        eventManager: CoroutineEventManager,
        eventListeners: Array<EventListener>,
        asyncEventListeners: Array<CoroutineEventListener>,
    ): JDABuilder = runBlocking {
        val bot = JDABuilder.createDefault(props.token)

        val listenerNames =
            arrayOf(*eventListeners, *asyncEventListeners).map { event ->
                event.javaClass.simpleName
            }
        logger.debug(
            StringMapMessage()
                .with("message", "Registering Event Listeners: $listenerNames")
                .with("event_listeners", eventListeners.map { it.javaClass.simpleName }.toString())
                .with(
                    "async_event_listeners",
                    asyncEventListeners.map { it.javaClass.simpleName }.toString(),
                )
        )
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
