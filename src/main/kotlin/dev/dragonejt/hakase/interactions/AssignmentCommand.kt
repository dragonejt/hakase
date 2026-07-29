import dev.dragonejt.hakase.interactions.ApplicationCommand
import dev.dragonejt.hakase.telemetry.LogBase
import dev.minn.jda.ktx.events.CoroutineEventListener
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.Commands

class AssignmentCommand : ApplicationCommand, CoroutineEventListener, LogBase() {
    override fun command() =
        Commands.slash("assignments", "List all Assigments")
            .addOption(OptionType.STRING, "Assignment ID", "ID of the assignment to view")

    override suspend fun onEvent(event: GenericEvent) {
        if (event !is SlashCommandInteractionEvent || event.fullCommandName != command().name)
            return
        event.deferReply().queue()

        val assignmentId = event.getOption("Assignment ID")?.asString
        when (assignmentId) {
            null -> listAll(event)
            else -> viewAssignment(event, assignmentId)
        }
    }

    @Suppress("UnusedPrivateMember", "UNUSED_PARAMETER")
    private suspend fun listAll(event: SlashCommandInteractionEvent) {
        // List all functions
    }

    @Suppress("UnusedPrivateMember", "UNUSED_PARAMETER")
    private suspend fun viewAssignment(event: SlashCommandInteractionEvent, assignmentId: String) {
        // Implement viewing a single assignment
    }
}
