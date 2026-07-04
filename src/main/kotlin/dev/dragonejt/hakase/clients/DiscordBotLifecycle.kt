package dev.dragonejt.hakase.clients

import dev.dragonejt.hakase.interactions.SlashCommand
import io.github.oshai.kotlinlogging.KotlinLogging
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.cancel
import net.dv8tion.jda.api.JDA
import net.dv8tion.jda.api.JDABuilder
import org.springframework.context.SmartLifecycle
import org.springframework.stereotype.Component

@Component
class DiscordBotLifecycle(
    private val scope: CoroutineScope,
    private val botConfig: JDABuilder,
    private val commands: List<SlashCommand>,
) : SmartLifecycle {
    private val log = KotlinLogging.logger {}
    private var bot: JDA? = null

    override fun start() {
        log.info { "Starting Discord Bot..." }

        bot = botConfig.build()
        bot!!.updateCommands().addCommands(commands.map { command -> command.command() }).queue()
        bot!!.awaitReady()
    }

    override fun stop() {
        stop {}
    }

    override fun stop(callback: Runnable) {
        log.info { "Stopping Discord Bot..." }
        bot!!.shutdown()
        bot!!.awaitShutdown()
        scope.cancel()
        callback.run()
    }

    override fun isRunning(): Boolean = bot?.status == JDA.Status.CONNECTED

    override fun isAutoStartup() = true
}
