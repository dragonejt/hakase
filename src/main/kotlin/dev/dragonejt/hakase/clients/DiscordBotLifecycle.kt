package dev.dragonejt.hakase.clients

import dev.kord.core.Kord
import io.github.oshai.kotlinlogging.KotlinLogging
import java.util.concurrent.atomic.AtomicBoolean
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.cancel
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking
import org.springframework.context.SmartLifecycle
import org.springframework.stereotype.Component

@Component
class DiscordBotLifecycle(private val bot: Kord) : SmartLifecycle {
    private val log = KotlinLogging.logger {}
    private val scope: CoroutineScope = CoroutineScope(SupervisorJob() + Dispatchers.IO)
    private var running = AtomicBoolean(false)

    override fun start() {
        log.info { "Starting Discord Bot..." }
        if (!running.compareAndSet(false, true)) return

        scope.launch {
            bot.login()
            running.set(false)
        }
    }

    override fun stop() {
        stop {}
    }

    override fun stop(callback: Runnable) {
        log.info { "Stopping Discord Bot..." }
        runBlocking {
            bot.logout()
            scope.cancel()
            running.set(false)
            callback.run()
        }
    }

    override fun isRunning(): Boolean = running.get()

    override fun isAutoStartup() = true
}
