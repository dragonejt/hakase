package dev.dragonejt.hakase.events

import dev.minn.jda.ktx.events.CoroutineEventListener
import io.github.oshai.kotlinlogging.KotlinLogging
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.session.ReadyEvent
import org.springframework.stereotype.Service

@Service
class ReadyHandler(private val tracer: Tracer) : CoroutineEventListener {
    private val log = KotlinLogging.logger {}

    override suspend fun onEvent(event: GenericEvent) {
        if (event !is ReadyEvent) return

        val span = tracer.spanBuilder("events.ready").setSpanKind(SpanKind.SERVER).startSpan()
        val scope = span.makeCurrent()

        log.info { "Logged in as ${event.jda.selfUser.name}!" }

        scope.close()
        span.end()
    }
}
