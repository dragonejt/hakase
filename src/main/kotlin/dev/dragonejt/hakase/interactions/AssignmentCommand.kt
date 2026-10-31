package dev.dragonejt.hakase.interactions

import dev.dragonejt.hakase.data.Assignment
import dev.dragonejt.hakase.data.AssignmentRepository
import dev.dragonejt.hakase.data.CourseRepository
import dev.dragonejt.hakase.telemetry.LogBase
import dev.minn.jda.ktx.events.CoroutineEventListener
import kotlin.time.toKotlinInstant
import net.dv8tion.jda.api.components.label.Label
import net.dv8tion.jda.api.components.textinput.TextInput
import net.dv8tion.jda.api.components.textinput.TextInputStyle
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.interaction.ModalInteractionEvent
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.Commands
import net.dv8tion.jda.api.interactions.commands.build.SubcommandData
import net.dv8tion.jda.api.modals.Modal
import org.natty.Parser
import org.springframework.stereotype.Service

@Service
class AssignmentCommand(
    private val courses: CourseRepository,
    private val assignments: AssignmentRepository,
) : ApplicationCommand, CoroutineEventListener, LogBase() {
    override fun command() =
        Commands.slash("assignments", "Manage assignments")
            .addSubcommands(
                SubcommandData("view", "View an assignment or list all assignments")
                    .addOption(
                        OptionType.STRING,
                        "assignment_id",
                        "ID of the assignment to view",
                        false,
                    ),
                SubcommandData("create", "Create a new assignment"),
            )

    override suspend fun onEvent(event: GenericEvent) {
        if (event is ModalInteractionEvent && event.modalId == "create_assignment_modal") {
            handleCreateModal(event)
            return
        }

        if (event !is SlashCommandInteractionEvent || event.name != command().name) return

        when (event.subcommandName) {
            "view" -> {
                val assignmentId = event.getOption("assignment_id")?.asString
                if (assignmentId == null) {
                    listAll(event)
                } else {
                    viewAssignment(event, assignmentId.toLong())
                }
            }
            "create" -> createAssignment(event)
        }
    }

    @Suppress("MagicNumber")
    private suspend fun listAll(event: SlashCommandInteractionEvent) {
        event.deferReply().queue()

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

        courseAssignments.take(25).forEach { assignment ->
            val value =
                "Due: <t:${assignment.due}:F>\nStatus: ${assignment.status}\nID: `${assignment.id}`"
            embedBuilder.addField(assignment.name, value, false)
        }

        event.hook.sendMessageEmbeds(embedBuilder.build()).queue()
    }

    private fun createAssignment(event: SlashCommandInteractionEvent) {
        val title =
            TextInput.create("name", TextInputStyle.SHORT)
                .setPlaceholder("e.g. Chapter 2 Homework")
                .build()
        val dueDate =
            TextInput.create("due_date", TextInputStyle.SHORT)
                .setPlaceholder("e.g. tomorrow at 5pm")
                .build()
        val url =
            TextInput.create("url", TextInputStyle.SHORT)
                .setPlaceholder("Must start with either http:// or https://")
                .setRequired(false)
                .build()

        val modal =
            Modal.create("create_assignment_modal", "Create Assignment")
                .addComponents(
                    Label.of("Assignment Name", title),
                    Label.of("Due Date", dueDate),
                    Label.of("URL", url),
                )
                .build()

        event.replyModal(modal).queue()
    }

    @Suppress("ReturnCount")
    private suspend fun handleCreateModal(event: ModalInteractionEvent) {
        event.deferReply().queue()

        val name = event.getValue("name")?.asString ?: "Untitled"
        val url = event.getValue("url")?.asString?.ifBlank { null }

        if (url != null && !url.startsWith("http://") && !url.startsWith("https://")) {
            event.hook
                .sendMessage(
                    "Invalid URL format! \nPlease ensure it starts with `http://` or `https://`."
                )
                .queue()
            return
        }

        val courseId =
            event.guild?.id
                ?: run {
                    event.hook.sendMessage("This command can only be used inside a server!").queue()
                    return
                }
        val status = "Not Started"

        val dueDateInput = event.getValue("due_date")?.asString?.trim()
        val parsedDueDate =
            dueDateInput
                ?.takeIf { it.isNotEmpty() }
                ?.let { Parser().parse(it).firstOrNull()?.dates?.firstOrNull() }

        if (parsedDueDate == null) {
            event.hook
                .sendMessage(
                    "I couldn't understand that due date. Try something like `tomorrow at 5pm`."
                )
                .queue()
            return
        }

        val assignment =
            Assignment(
                null,
                courses.findById(courseId).get(),
                parsedDueDate.toInstant().toKotlinInstant(),
                name,
                status,
                url,
            )
        val savedAssignment = assignments.save(assignment)

        event.hook
            .sendMessage(
                "Assignment created successfully: **$name** \n(ID: `${savedAssignment.id}`)"
            )
            .queue()
    }

    private suspend fun viewAssignment(event: SlashCommandInteractionEvent, assignmentId: Long) {
        event.deferReply().queue()

        val assignmentOpt = assignments.findById(assignmentId)

        if (assignmentOpt.isEmpty) {
            event.hook.sendMessage("Assignment not found.").queue()
            return
        }

        val assignment = assignmentOpt.get()

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
