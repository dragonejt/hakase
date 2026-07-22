package dev.dragonejt.hakase.clients

import dev.dragonejt.hakase.interactions.ApplicationCommand
import net.dv8tion.jda.api.JDA
import net.dv8tion.jda.api.JDABuilder
import net.dv8tion.jda.api.interactions.commands.build.CommandData
import net.dv8tion.jda.api.interactions.commands.build.Commands
import net.dv8tion.jda.api.requests.restaction.CommandListUpdateAction
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertAll
import org.junit.jupiter.api.extension.ExtendWith
import org.mockito.Mock
import org.mockito.junit.jupiter.MockitoExtension
import org.mockito.kotlin.any
import org.mockito.kotlin.mock
import org.mockito.kotlin.verify
import org.mockito.kotlin.whenever

@ExtendWith(MockitoExtension::class)
class DiscordBotLifecycleTests {
    @Mock private lateinit var botConfig: JDABuilder
    @Mock private lateinit var bot: JDA
    @Mock private lateinit var command: ApplicationCommand
    @Mock private lateinit var updateCommands: CommandListUpdateAction

    private lateinit var underTest: DiscordBotLifecycle

    @BeforeEach
    fun setUp() {
        whenever { command.command() }.thenReturn(Commands.slash("hakase", "hakase"))
        whenever { botConfig.build() }.thenReturn(bot)
        whenever { bot.updateCommands() }.thenReturn(updateCommands)
        whenever { updateCommands.addCommands(any<Collection<CommandData>>()) }
            .thenReturn(mock<CommandListUpdateAction>())
        underTest = DiscordBotLifecycle(botConfig, listOf(command))
    }

    @Test
    @DisplayName("start() correctly initializes JDA bot")
    fun testStartInitializesBot() {
        underTest.start()

        assertAll(
            { verify(botConfig).build() },
            { verify(bot).awaitReady() },
            { assertThat(underTest.isRunning()).isTrue() },
        )
    }

    @Test
    @DisplayName("start() correctly registers app commands")
    fun testStartRegistersCommands() {
        underTest.start()

        assertAll(
            { verify(updateCommands).addCommands(any<Collection<CommandData>>()) },
            { verify(command).command() },
        )
    }

    @Test
    @DisplayName("stop() correctly stops discord bot")
    fun testStopStopsDiscordBot() {
        underTest.start()
        assertThat(underTest.isRunning()).isTrue()

        underTest.stop()

        assertAll(
            { verify(bot).shutdown() },
            { verify(bot).awaitShutdown() },
            { assertThat(underTest.isRunning()).isFalse() },
        )
    }

    @Test
    @DisplayName("stop() runs callback after stop")
    fun testStopRunsCallback() {
        underTest.start()
        assertThat(underTest.isRunning()).isTrue()

        val callback = mock<Runnable>()
        underTest.stop(callback)

        verify(callback).run()
    }
}
