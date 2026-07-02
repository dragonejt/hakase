package dev.dragonejt.hakase.interactions

import dev.minn.jda.ktx.coroutines.await
import dev.minn.jda.ktx.events.CoroutineEventListener
import io.github.oshai.kotlinlogging.KotlinLogging
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.Commands
import net.dv8tion.jda.api.interactions.commands.build.OptionData
import org.springframework.stereotype.Service

@Service
class HakaseCommand(private val tracer: Tracer, private val scope: CoroutineScope) :
    SlashCommand, CoroutineEventListener {
    private val log = KotlinLogging.logger {}

    override fun command() =
        Commands.slash("hakase", "hakase settings")
            .addOptions(
                OptionData(OptionType.STRING, "cmd", "subcommand to run")
                    .addChoice(
                        "config",
                        "hakase configuration",
                    )
                    .addChoice("rps", "rock-paper-scissors")
            )

    override suspend fun onEvent(event: GenericEvent) {
        if (event !is SlashCommandInteractionEvent || event.fullCommandName != command().name)
            return

        scope.launch {
            val span =
                tracer.spanBuilder("commands.hakase").setSpanKind(SpanKind.SERVER).startSpan()
            val scope = span.makeCurrent()
            log.info { "/hakase executed by: ${event.user.effectiveName}" }

            event.reply("hakase pong!").await()

            scope.close()
            span.end()
        }
    }
}
