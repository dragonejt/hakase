package dev.dragonejt.hakase.events

import dev.kord.core.Kord
import dev.kord.core.entity.effectiveName
import dev.kord.core.event.gateway.ReadyEvent
import dev.kord.core.on
import io.github.oshai.kotlinlogging.KotlinLogging
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import org.springframework.stereotype.Service

@Service
class ReadyHandler(private val tracer: Tracer) : EventHandler<ReadyEvent> {
    private val log = KotlinLogging.logger {}

    override fun register(bot: Kord) {
        bot.on<ReadyEvent> {
            val span = tracer.spanBuilder("events.ready").setSpanKind(SpanKind.SERVER).startSpan()
            val scope = span.makeCurrent()

            handleEvent(this)

            scope.close()
            span.end()
        }
    }

    override fun handleEvent(event: ReadyEvent) {
        log.info { "Logged in as ${event.self.effectiveName}!" }
    }
}
