package dev.dragonejt.hakase.events

import dev.dragonejt.hakase.telemetry.LogBase
import dev.minn.jda.ktx.events.CoroutineEventListener
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import net.dv8tion.jda.api.entities.Activity
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.session.ReadyEvent
import org.springframework.stereotype.Service

@Service
class ReadyHandler(private val tracer: Tracer) : CoroutineEventListener, LogBase() {

    override suspend fun onEvent(event: GenericEvent) {
        if (event !is ReadyEvent) return

        val span =
            tracer
                .spanBuilder("events.${this.javaClass.simpleName}")
                .setSpanKind(SpanKind.SERVER)
                .startSpan()
        val scope = span.makeCurrent()

        log.atInfo {
            message = "Logged in as ${event.jda.selfUser.name}!"
            payload = mapOf("bot_user" to event.jda.selfUser.name)
        }

        event.jda.presence.activity =
            Activity.of(
                Activity.ActivityType.CUSTOM_STATUS,
                "hakase.dragonejt.dev | in ${event.jda.guilds.size} courses",
            )

        scope.close()
        span.end()
    }
}
