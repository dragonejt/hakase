package dev.dragonejt.hakase.interactions

import dev.kord.core.Kord
import dev.kord.core.entity.interaction.Interaction
import dev.kord.core.event.interaction.GuildChatInputCommandInteractionCreateEvent
import dev.kord.core.on
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer

class HakaseCommand(private val tracer: Tracer) :
    InteractionHandler<GuildChatInputCommandInteractionCreateEvent> {

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

    override suspend fun handleInteraction(interaction: Interaction) {
        TODO("Not yet implemented")
    }
}
