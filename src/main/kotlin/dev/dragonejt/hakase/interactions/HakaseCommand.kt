package dev.dragonejt.hakase.interactions

import dev.dragonejt.hakase.clients.DiscordProperties
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
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.stereotype.Service

@Service
@EnableConfigurationProperties(DiscordProperties::class)
class HakaseCommand(
    private val tracer: Tracer,
    private val scope: CoroutineScope,
    private val props: DiscordProperties,
) : ApplicationCommand, CoroutineEventListener, LogBase() {

    enum class SubCommand(val cmd: String?) {
        CONFIG("config"),
        RPS("rps"),
        PONG(null),
    }

    override fun command() =
        Commands.slash("hakase", "hakase settings")
            .addOptions(
                OptionData(OptionType.STRING, "cmd", "subcommand to run")
                    .addChoice(
                        SubCommand.CONFIG.cmd!!,
                        SubCommand.CONFIG.cmd,
                    )
                    .addChoice(SubCommand.RPS.cmd!!, SubCommand.RPS.cmd)
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

        val cmdOption = SubCommand.entries.find { it.cmd == event.getOption("cmd")?.asString }
        when (cmdOption) {
            SubCommand.RPS -> rps(event)
            else -> default(event)
        }

        scope.close()
        span.end()
    }

    private suspend fun rps(event: SlashCommandInteractionEvent) {
        event.reply(props.rpsGifs.random()).await()
    }

    private suspend fun default(event: SlashCommandInteractionEvent) {
        event.reply("hakase pong!").await()
    }
}
