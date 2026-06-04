package dev.dragonejt.hakase.events

import dev.kord.core.Kord
import dev.kord.core.entity.effectiveName
import dev.kord.core.event.gateway.ReadyEvent
import dev.kord.core.on
import io.github.oshai.kotlinlogging.KotlinLogging
import org.springframework.stereotype.Service

@Service
class ReadyHandler : EventHandler<ReadyEvent> {

  private val log = KotlinLogging.logger {}

  override fun register(bot: Kord) {
    bot.on<ReadyEvent> { handleEvent(this) }
  }

  override fun handleEvent(event: ReadyEvent) {
    log.info { "Logged in as ${event.self.effectiveName}!" }
  }
}
