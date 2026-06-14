package dev.dragonejt.hakase.clients

import dev.dragonejt.hakase.events.EventHandler
import dev.kord.core.Kord
import dev.kord.core.event.gateway.GatewayEvent
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import org.mockito.Mock
import org.mockito.Mockito.mock
import org.mockito.junit.jupiter.MockitoExtension
import org.mockito.kotlin.verify

@ExtendWith(MockitoExtension::class)
class DiscordBotConfigurationTests {

    @Mock private lateinit var bot: Kord

    private val properties = DiscordProperties(token = "DISCORD_BOT_TOKEN")

    private lateinit var eventHandlers: List<EventHandler<out GatewayEvent>>

    private lateinit var underTest: DiscordBotConfiguration

    @BeforeEach
    fun setUp() {
        eventHandlers = listOf(mock<EventHandler<out GatewayEvent>>(EventHandler::class.java))

        underTest = DiscordBotConfiguration(this::mockKordFactory)
    }

    @Test
    @DisplayName("Proper event handlers are loaded")
    fun testEventHandlersLoaded() {
        underTest.discordBot(properties, eventHandlers)

        eventHandlers.forEach { handler -> verify(handler).register(bot) }
    }

    private fun mockKordFactory(token: String): Kord = this.bot
}
