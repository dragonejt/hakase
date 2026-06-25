package dev.dragonejt.hakase.interactions

import dev.kord.core.Kord
import dev.kord.core.behavior.interaction.response.respond
import dev.kord.core.entity.interaction.GuildChatInputCommandInteraction
import dev.kord.core.event.interaction.GuildChatInputCommandInteractionCreateEvent
import dev.kord.core.on
import io.github.oshai.kotlinlogging.KotlinLogging
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import org.springframework.stereotype.Service

@Service
class HakaseCommand(private val tracer: Tracer) :
    InteractionHandler<GuildChatInputCommandInteraction> {
    private val log = KotlinLogging.logger {}

    override suspend fun register(bot: Kord) {
        bot.createGlobalChatInputCommand("hakase", "hakase settings")

        bot.on<GuildChatInputCommandInteractionCreateEvent> {
            val span = tracer.spanBuilder("events.ready").setSpanKind(SpanKind.SERVER).startSpan()
            val scope = span.makeCurrent()
            handleInteraction(interaction)
            scope.close()
            span.end()
        }
    }

    override suspend fun handleInteraction(interaction: GuildChatInputCommandInteraction) {
        log.info { "/hakase executed by: ${interaction.user.username}" }
        val response = interaction.deferPublicResponse()
        response.respond { content = "hello" }
    }
}
