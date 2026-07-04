package dev.dragonejt.hakase.clients

import dev.dragonejt.hakase.interactions.ApplicationCommand
import dev.dragonejt.hakase.telemetry.LogBase
import java.util.concurrent.atomic.AtomicBoolean
import net.dv8tion.jda.api.JDA
import net.dv8tion.jda.api.JDABuilder
import org.springframework.context.SmartLifecycle
import org.springframework.stereotype.Component

@Component
class DiscordBotLifecycle(
    private val botConfig: JDABuilder,
    private val commands: List<ApplicationCommand>,
) : SmartLifecycle, LogBase() {
    private lateinit var bot: JDA
    private var running = AtomicBoolean(false)

    override fun start() {
        if (isRunning()) return

        log.info { "Starting Discord Bot..." }
        bot = botConfig.build()
        log.atDebug {
            message =
                "Registering Application Commands: " +
                    commands.map { command -> command.javaClass.simpleName }
            payload = mapOf("commands" to commands)
        }
        bot.updateCommands().addCommands(commands.map { command -> command.command() }).queue()
        bot.awaitReady()
        running.set(true)
    }

    override fun stop() {
        stop {}
    }

    override fun stop(callback: Runnable) {
        if (!isRunning()) return

        log.info { "Stopping Discord Bot..." }
        bot.shutdown()
        bot.awaitShutdown()
        running.set(false)
        callback.run()
    }

    override fun isRunning(): Boolean = running.get()

    override fun isAutoStartup() = true
}
