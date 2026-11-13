package dev.dragonejt.hakase.interactions

import dev.dragonejt.hakase.data.Assignment
import dev.dragonejt.hakase.data.AssignmentRepository
import dev.dragonejt.hakase.data.CourseRepository
import dev.dragonejt.hakase.displays.Display
import dev.minn.jda.ktx.events.CoroutineEventListener
import io.opentelemetry.api.trace.SpanKind
import io.opentelemetry.api.trace.Tracer
import net.dv8tion.jda.api.events.GenericEvent
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.Commands
import org.apache.logging.log4j.kotlin.Logging
import org.apache.logging.log4j.message.StringMapMessage
import org.springframework.stereotype.Service

@Service
class AssignmentCommand(
    private val tracer: Tracer,
    private val courses: CourseRepository,
    private val assignments: AssignmentRepository,
    private val assignmentsDisplay: Display<List<Assignment>>,
    private val assignmentDisplay: Display<Assignment>,
) : ApplicationCommand, CoroutineEventListener {

    companion object : Logging

    override fun command() =
        Commands.slash("assignments", "manage assignments")
            .addOption(
                OptionType.STRING,
                "id",
                "retrieves assignment with this id",
                false,
            )

    override suspend fun onEvent(event: GenericEvent) {
        val span =
            tracer
                .spanBuilder("events.${this.javaClass.simpleName}")
                .setSpanKind(SpanKind.SERVER)
                .startSpan()
        val scope = span.makeCurrent()

        if (event !is SlashCommandInteractionEvent || event.name != command().name) return
        event.deferReply().queue()

        val assignmentId = event.getOption("id")
        logger.info(
            StringMapMessage()
                .with("message", "/assignments has been executed")
                .with("assignment_id", assignmentId?.asLong)
                .with("course_id", event.guild?.id)
                .with("user", event.user.name)
        )

        if (assignmentId == null) {
            listAll(event)
        } else {
            viewAssignment(event, assignmentId.asLong)
        }

        scope.close()
        span.end()
    }

    private suspend fun listAll(event: SlashCommandInteractionEvent) {
        val courseId =
            event.guild?.id
                ?: run {
                    event.hook.sendMessage("This command can only be used inside a server!").queue()
                    return
                }
        val course = courses.findById(courseId).get()
        event.hook
            .sendMessageEmbeds(assignmentsDisplay.embed(course.assignments))
            .addComponents(assignmentsDisplay.component(course.assignments))
            .queue()
    }

    private suspend fun viewAssignment(event: SlashCommandInteractionEvent, assignmentId: Long) {
        val assignment = assignments.findById(assignmentId).orElse(null)
        assignment
            ?: run {
                event.hook.sendMessage("Assignment not found.").queue()
                return
            }

        event.hook
            .sendMessageEmbeds(assignmentDisplay.embed(assignment))
            .addComponents(assignmentDisplay.component(assignment))
            .queue()
    }
}
