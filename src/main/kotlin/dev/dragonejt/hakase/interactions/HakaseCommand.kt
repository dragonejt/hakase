package dev.dragonejt.hakase.interactions

import dev.dragonejt.hakase.telemetry.LogBase
import dev.minn.jda.ktx.coroutines.await
import dev.minn.jda.ktx.events.CoroutineEventListener
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import kotlinx.coroutines.CoroutineScope
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.Commands
import net.dv8tion.jda.api.interactions.commands.build.OptionData
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Service

@Service
class HakaseCommand(private val tracer: Tracer, private val scope: CoroutineScope) :
    ApplicationCommand, CoroutineEventListener, LogBase() {

    @Value("\${discord-bot.rps-gifs}") private lateinit var rpsGifs: List<String>

    override fun command() =
        Commands.slash("hakase", "hakase settings")
            .addOptions(
                OptionData(OptionType.STRING, "cmd", "subcommand to run")
                    .addChoice(
                        "config",
                        "hakase configuration",
                    )
                    .addChoice("rps", "rps")
            )

    override suspend fun onEvent(event: GenericEvent) {
        if (event !is SlashCommandInteractionEvent || event.fullCommandName != command().name)
            return

        val span =
            tracer
                .spanBuilder("commands.${this.javaClass.simpleName}")
                .setSpanKind(SpanKind.SERVER)
                .startSpan()
        val scope = span.makeCurrent()
        log.atInfo {
            message = "/hakase executed by: ${event.user.effectiveName}"
            payload = mapOf("username" to event.user.effectiveName)
        }

        val cmdOption = event.getOption("cmd")?.asString
        when (cmdOption) {
            "rps" -> event.reply(rps()).await()
            else -> event.reply("hakase pong!").await()
        }

        scope.close()
        span.end()
    }

    private fun rps(): String {
        return rpsGifs.random()
    }
}
