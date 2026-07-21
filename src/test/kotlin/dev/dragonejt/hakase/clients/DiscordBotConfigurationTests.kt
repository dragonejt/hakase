package dev.dragonejt.hakase.clients

import dev.minn.jda.ktx.events.CoroutineEventListener
import dev.minn.jda.ktx.events.CoroutineEventManager
import kotlin.random.Random
import net.dv8tion.jda.api.JDABuilder
import net.dv8tion.jda.api.hooks.EventListener
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertAll
import org.junit.jupiter.api.extension.ExtendWith
import org.mockito.Mock
import org.mockito.Mockito.mock
import org.mockito.Mockito.mockStatic
import org.mockito.junit.jupiter.MockitoExtension
import org.mockito.kotlin.any
import org.mockito.kotlin.verify
import org.mockito.kotlin.whenever

@ExtendWith(MockitoExtension::class)
class DiscordBotConfigurationTests {

    @Mock private lateinit var bot: JDABuilder

    private val properties =
        DiscordProperties(
            token = "DISCORD_BOT_TOKEN",
            rpsGifs =
                listOf("https://klipy.com/gifs/the-amazing-world-of-gumball-gumball-and-darwin"),
        )

    @Mock private lateinit var eventManager: CoroutineEventManager
    private lateinit var eventListeners: Array<EventListener>

    private lateinit var asyncEventListeners: Array<CoroutineEventListener>

    private lateinit var underTest: DiscordBotConfiguration

    @BeforeEach
    fun setUp() {
        eventListeners =
            Array(Random.nextInt(10)) {
                mock<EventListener>(EventListener::class.java)
            }
        asyncEventListeners =
            Array(Random.nextInt(10)) {
                mock<CoroutineEventListener>(CoroutineEventListener::class.java)
            }

        underTest = DiscordBotConfiguration()
    }

    @Test
    @DisplayName("Proper event handlers and interaction handlers are loaded")
    fun testHandlersLoaded() {
        mockStatic(JDABuilder::class.java).use { botBuilder ->
            botBuilder.whenever { JDABuilder.createDefault(any<String>()) }.thenReturn(bot)

            underTest.discordBot(properties, eventManager, eventListeners, asyncEventListeners)

            assertAll(
                { verify(bot).addEventListeners(*eventListeners) },
                { verify(bot).addEventListeners(*asyncEventListeners) },
                { verify(bot).setEventManager(eventManager) },
            )
        }
    }
}
