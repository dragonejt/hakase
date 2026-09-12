package dev.dragonejt.hakase.interactions

import dev.dragonejt.hakase.datamodels.Assignment
import dev.dragonejt.hakase.osdk.AssignmentRepository
import dev.dragonejt.hakase.telemetry.LogBase
import dev.minn.jda.ktx.events.CoroutineEventListener
import java.time.LocalDateTime
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import net.dv8tion.jda.api.components.label.Label
import net.dv8tion.jda.api.components.textinput.TextInput
import net.dv8tion.jda.api.components.textinput.TextInputStyle
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.interaction.ModalInteractionEvent
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.Commands
import net.dv8tion.jda.api.modals.Modal
import org.springframework.stereotype.Service

@Service
class AssignmentCommand(private val assignments: AssignmentRepository) :
    ApplicationCommand, CoroutineEventListener, LogBase() {
    override fun command() =
        Commands.slash("assignments", "List all Assigments")
            .addOption(OptionType.STRING, "assignment_id", "ID of the assignment to view")

    override suspend fun onEvent(event: GenericEvent) {
        if (event is ModalInteractionEvent && event.modalId == "create_assignment_modal") {
            handleCreateModal(event)
            return
        }

        if (event !is SlashCommandInteractionEvent || event.fullCommandName != command().name)
            return

        when (val assignmentId = event.getOption("assignment_id")?.asString) {
            null -> listAll(event)
            "create" -> createAssignment(event)
            else -> viewAssignment(event, assignmentId)
        }
    }

    @Suppress("UnusedPrivateMember", "UNUSED_PARAMETER")
    private suspend fun listAll(event: SlashCommandInteractionEvent) {
        event.deferReply().queue()
    }

    private fun createAssignment(event: SlashCommandInteractionEvent) {
        val title =
            TextInput.create("name", TextInputStyle.SHORT)
                .setPlaceholder("e.g. Chapter 2 Homework")
                .build()
        val dueDate =
            TextInput.create("due_date", TextInputStyle.SHORT)
                .setPlaceholder("MM-DD-YYYY HH:mm (24hr format)")
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
        val dueDateString = event.getValue("due_date")?.asString
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

        val formatter = DateTimeFormatter.ofPattern("MM-dd-yyyy HH:mm")
        val dueDate =
            try {
                LocalDateTime.parse(dueDateString, formatter)
                    .atZone(ZoneId.systemDefault())
                    .toOffsetDateTime()
            } catch (e: java.time.format.DateTimeParseException) {
                event.hook
                    .sendMessage(
                        "Invalid date format! \nPlease use MM-dd-yyyy HH:mm (e.g. `09-30-2026 23:59`)."
                    )
                    .queue()
                return
            }
        val assignment = Assignment("", courseId, dueDate, name, status, url)
        val savedAssignment = assignments.save(assignment)

        event.hook
            .sendMessage(
                "Assignment created successfully: **$name** \n(ID: `${savedAssignment.id}`)"
            )
            .queue()
    }

    private suspend fun viewAssignment(event: SlashCommandInteractionEvent, assignmentId: String) {
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
