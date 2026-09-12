package dev.dragonejt.hakase.events

import dev.dragonejt.hakase.osdk.CourseRepository
import dev.dragonejt.hakase.telemetry.LogBase
import dev.minn.jda.ktx.events.CoroutineEventListener
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import net.dv8tion.jda.api.entities.Activity
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.guild.GuildLeaveEvent
import org.springframework.stereotype.Service

@Service
class GuildLeaveHandler(private val tracer: Tracer, private val courses: CourseRepository) :
    CoroutineEventListener, LogBase() {
    override suspend fun onEvent(event: GenericEvent) {
        if (event !is GuildLeaveEvent) return

        val span =
            tracer
                .spanBuilder("events.${this.javaClass.simpleName}")
                .setSpanKind(SpanKind.SERVER)
                .startSpan()
        val scope = span.makeCurrent()

        log.atInfo {
            message = "Bot left guild ${event.guild.name} (${event.guild.id})"
            payload = mapOf("guild_name" to event.guild.name, "guild_id" to event.guild.id)
        }
        courses.deleteById(event.guild.id)
        event.jda.presence.activity =
            Activity.of(
                Activity.ActivityType.CUSTOM_STATUS,
                "hakase.dragonejt.dev | in ${event.jda.guilds.size} courses",
            )

        scope.close()
        span.end()
    }
}
