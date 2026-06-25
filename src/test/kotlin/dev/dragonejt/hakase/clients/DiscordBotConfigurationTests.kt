package dev.dragonejt.hakase.clients

import dev.dragonejt.hakase.events.EventHandler
import dev.dragonejt.hakase.interactions.InteractionHandler
import dev.kord.core.Kord
import dev.kord.core.entity.interaction.Interaction
import dev.kord.core.event.gateway.GatewayEvent
import kotlin.random.Random
import org.assertj.core.api.Assertions.assertThat
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

    private lateinit var interactionHandlers: List<InteractionHandler<out Interaction>>

    private lateinit var underTest: DiscordBotConfiguration

    @BeforeEach
    fun setUp() {
        eventHandlers =
            List(Random.nextInt(10)) {
                mock<EventHandler<out GatewayEvent>>(EventHandler::class.java)
            }
        interactionHandlers =
            List(Random.nextInt(10)) {
                mock<InteractionHandler<out Interaction>>(InteractionHandler::class.java)
            }

        underTest = DiscordBotConfiguration(this::mockKordFactory)
    }

    @Test
    @DisplayName("Proper event handlers and interaction handlers are loaded")
    suspend fun testHandlersLoaded() {
        underTest.discordBot(properties, eventHandlers, interactionHandlers)

        eventHandlers.forEach { handler -> verify(handler).register(bot) }
        interactionHandlers.forEach { handler -> verify(handler).register(bot) }
    }

    private fun mockKordFactory(token: String): Kord {
        assertThat(token).isEqualTo(properties.token)
        return bot
    }
}
