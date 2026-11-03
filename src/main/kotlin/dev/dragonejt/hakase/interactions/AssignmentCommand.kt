package dev.dragonejt.hakase.interactions

import dev.dragonejt.hakase.data.AssignmentRepository
import dev.dragonejt.hakase.data.CourseRepository
import dev.dragonejt.hakase.telemetry.LogBase
import dev.minn.jda.ktx.events.CoroutineEventListener
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.Commands
import org.springframework.stereotype.Service

@Service
class AssignmentCommand(
    private val courses: CourseRepository,
    private val assignments: AssignmentRepository,
) : ApplicationCommand, CoroutineEventListener, LogBase() {
    override fun command() =
        Commands.slash("assignments", "Manage assignments")
            .addOption(
                OptionType.STRING,
                "id",
                "retrieves assignment with this id",
                false,
            )

    override suspend fun onEvent(event: GenericEvent) {
        if (event !is SlashCommandInteractionEvent || event.name != command().name) return

        val assignmentId = event.getOption("id")
        if (assignmentId == null) {
            listAll(event)
        } else {
            viewAssignment(event, assignmentId.asLong)
        }
    }

    private suspend fun listAll(event: SlashCommandInteractionEvent) {
        val courseId =
            event.guild?.id
                ?: run {
                    event.hook.sendMessage("This command can only be used inside a server!").queue()
                    return
                }
        val course = courses.findById(courseId).get()

        val courseAssignments = course.assignments

        if (courseAssignments.isEmpty()) {
            event.hook.sendMessage("No assignments found for this course!").queue()
            return
        }

        val embedBuilder =
            net.dv8tion.jda.api.EmbedBuilder().setTitle("Assignments").setColor(0x5865F2)

        courseAssignments.forEach { assignment ->
            val value =
                "Due: <t:${assignment.due}:F>\nStatus: ${assignment.status}\nID: `${assignment.id}`"
            embedBuilder.addField(assignment.name, value, false)
        }

        event.hook.sendMessageEmbeds(embedBuilder.build()).queue()
    }

    private suspend fun viewAssignment(event: SlashCommandInteractionEvent, assignmentId: Long) {
        val assignment = assignments.findById(assignmentId).orElse(null)
        assignment
            ?: run {
                event.hook.sendMessage("Assignment not found.").queue()
                return
            }

        val embed =
            net.dv8tion.jda.api
                .EmbedBuilder()
                .setTitle(assignment.name, assignment.url)
                .addField("Course ID", assignment.course.courseId, true)
                .addField("Status", assignment.status, true)
                .addField("Due Date", "<t:${assignment.due}:F>", false)
                .setFooter("ID: ${assignment.id}")
                .setColor(0x5865F2)
                .build()

        event.hook.sendMessageEmbeds(embed).queue()
    }
}
