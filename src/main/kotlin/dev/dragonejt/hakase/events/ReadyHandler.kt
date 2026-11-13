package dev.dragonejt.hakase.events

import dev.minn.jda.ktx.events.CoroutineEventListener
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import net.dv8tion.jda.api.entities.Activity
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.session.ReadyEvent
import org.apache.logging.log4j.kotlin.Logging
import org.apache.logging.log4j.message.StringMapMessage
import org.springframework.stereotype.Service

@Service
class ReadyHandler(private val tracer: Tracer) : CoroutineEventListener {
    companion object : Logging

    override suspend fun onEvent(event: GenericEvent) {
        if (event !is ReadyEvent) return

        val span =
            tracer
                .spanBuilder("events.${this.javaClass.simpleName}")
                .setSpanKind(SpanKind.SERVER)
                .startSpan()
        val scope = span.makeCurrent()

        logger.info(
            StringMapMessage()
                .with("message", "Logged in as ${event.jda.selfUser.name}!")
                .with("bot_user", event.jda.selfUser.name)
        )

        event.jda.presence.activity =
            Activity.of(
                Activity.ActivityType.CUSTOM_STATUS,
                "hakase.dragonejt.dev | in ${event.jda.guilds.size} courses",
            )

        scope.close()
        span.end()
    }
}
