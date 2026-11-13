package dev.dragonejt.hakase.clients

import dev.dragonejt.hakase.interactions.ApplicationCommand
import java.util.concurrent.atomic.AtomicBoolean
import net.dv8tion.jda.api.JDA
import net.dv8tion.jda.api.JDABuilder
import org.apache.logging.log4j.kotlin.Logging
import org.apache.logging.log4j.message.StringMapMessage
import org.springframework.context.SmartLifecycle
import org.springframework.stereotype.Component

@Component
class DiscordBotLifecycle(
    private val botConfig: JDABuilder,
    private val commands: List<ApplicationCommand>,
) : SmartLifecycle {
    companion object : Logging

    private lateinit var bot: JDA
    private var running = AtomicBoolean(false)

    override fun start() {
        if (isRunning()) return

        logger.info { "Starting Discord Bot..." }
        bot = botConfig.build()
        val commandNames = commands.map { command -> command.javaClass.simpleName }
        logger.debug(
            StringMapMessage()
                .with("message", "Registering Application Commands: $commandNames")
                .with("commands", commandNames.toString())
        )
        bot.updateCommands().addCommands(commands.map { command -> command.command() }).queue()
        bot.awaitReady()
        running.set(true)
    }

    override fun stop() {
        stop {}
    }

    override fun stop(callback: Runnable) {
        if (!isRunning()) return

        logger.info { "Stopping Discord Bot..." }
        bot.shutdown()
        bot.awaitShutdown()
        running.set(false)
        callback.run()
    }

    override fun isRunning(): Boolean = running.get()

    override fun isAutoStartup() = true
}
