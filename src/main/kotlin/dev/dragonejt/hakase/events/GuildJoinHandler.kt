package dev.dragonejt.hakase.events

import dev.dragonejt.hakase.data.Course
import dev.dragonejt.hakase.data.CourseRepository
import dev.minn.jda.ktx.events.CoroutineEventListener
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import net.dv8tion.jda.api.entities.Activity
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.guild.GuildJoinEvent
import org.apache.logging.log4j.kotlin.Logging
import org.apache.logging.log4j.message.StringMapMessage
import org.springframework.stereotype.Service

@Service
class GuildJoinHandler(private val tracer: Tracer, private val courses: CourseRepository) :
    CoroutineEventListener {
    companion object : Logging

    override suspend fun onEvent(event: GenericEvent) {
        if (event !is GuildJoinEvent) return

        val span =
            tracer
                .spanBuilder("events.${this.javaClass.simpleName}")
                .setSpanKind(SpanKind.SERVER)
                .startSpan()
        val scope = span.makeCurrent()

        logger.info(
            StringMapMessage()
                .with("message", "Bot joined guild ${event.guild.name} (${event.guild.id})")
                .with("guild_name", event.guild.name)
                .with("guild_id", event.guild.id)
        )

        val course = Course(event.guild.id, event.guild.defaultChannel!!.id, "")
        courses.save(course)

        event.jda.presence.activity =
            Activity.of(
                Activity.ActivityType.CUSTOM_STATUS,
                "hakase.dragonejt.dev | in ${event.jda.guilds.size} courses",
            )

        scope.close()
        span.end()
    }
}
