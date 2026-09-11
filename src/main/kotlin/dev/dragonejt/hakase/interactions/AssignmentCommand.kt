package dev.dragonejt.hakase.interactions

import dev.dragonejt.hakase.osdk.AssignmentRepository
import dev.dragonejt.hakase.telemetry.LogBase
import dev.minn.jda.ktx.events.CoroutineEventListener
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.Commands
import org.springframework.stereotype.Service

@Service
class AssignmentCommand(private val assignments: AssignmentRepository) :
    ApplicationCommand, CoroutineEventListener, LogBase() {
    override fun command() =
        Commands.slash("assignments", "List all Assigments")
            .addOption(OptionType.STRING, "assignment_id", "ID of the assignment to view")

    override suspend fun onEvent(event: GenericEvent) {
        if (event !is SlashCommandInteractionEvent || event.fullCommandName != command().name)
            return
        event.deferReply().queue()

        when (val assignmentId = event.getOption("assignment_id")?.asString) {
            null -> listAll(event)
            else -> viewAssignment(event, assignmentId)
        }
    }

    @Suppress("UnusedPrivateMember", "UNUSED_PARAMETER")
    private suspend fun listAll(event: SlashCommandInteractionEvent) {
        // List all functions
    }

    private suspend fun viewAssignment(event: SlashCommandInteractionEvent, assignmentId: String) {
        val assignmentOpt = assignments.findById(assignmentId)

        if (assignmentOpt.isEmpty) {
            event.hook.sendMessage("Assignment not found.").queue()
            return
        }

        val assignment = assignmentOpt.get()

        val embed =
            net.dv8tion.jda.api
                .EmbedBuilder()
                .setTitle(assignment.name, assignment.url?.ifBlank { null })
                .addField("Course ID", assignment.courseID, true)
                .addField("Status", assignment.status, true)
                .addField("Due Date", "<t:${assignment.dueDate.toEpochSecond()}:F>", false)
                .setFooter("ID: ${assignment.id}")
                .setColor(0x5865F2)
                .build()

        event.hook.sendMessageEmbeds(embed).queue()
    }
}
