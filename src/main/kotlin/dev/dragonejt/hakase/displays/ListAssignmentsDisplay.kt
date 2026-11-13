package dev.dragonejt.hakase.displays

import dev.dragonejt.hakase.data.Assignment
import dev.minn.jda.ktx.messages.EmbedBuilder
import java.time.format.DateTimeFormatter
import net.dv8tion.jda.api.components.MessageTopLevelComponent
import net.dv8tion.jda.api.components.actionrow.ActionRow
import net.dv8tion.jda.api.components.buttons.Button
import net.dv8tion.jda.api.entities.MessageEmbed
import org.springframework.stereotype.Component

@Component
class ListAssignmentsDisplay : Display<MutableList<Assignment>> {

    override fun embed(toDisplay: MutableList<Assignment>): MessageEmbed {
        val embed = EmbedBuilder()
        embed.title = "Assignments"
        embed.description = "${toDisplay.size} Assignments in Course"

        toDisplay.forEach { assignment ->
            embed.field("ID: ${assignment.id}", assignment.name, true)
            embed.field(
                "Due",
                DateTimeFormatter.RFC_1123_DATE_TIME.format(assignment.due),
                false,
            )
        }

        return embed.build()
    }

    override fun component(toDisplay: MutableList<Assignment>): MessageTopLevelComponent {
        return ActionRow.of(Button.primary("add_assignment", "➕ Add"))
    }
}
