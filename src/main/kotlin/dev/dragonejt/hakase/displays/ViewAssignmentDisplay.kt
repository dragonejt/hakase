package dev.dragonejt.hakase.displays

import dev.dragonejt.hakase.data.Assignment
import dev.minn.jda.ktx.interactions.components.ActionRow
import dev.minn.jda.ktx.messages.EmbedBuilder
import net.dv8tion.jda.api.components.MessageTopLevelComponent
import net.dv8tion.jda.api.components.buttons.Button
import net.dv8tion.jda.api.entities.MessageEmbed
import org.springframework.stereotype.Component

@Component
class ViewAssignmentDisplay : Display<Assignment> {
    override fun embed(toDisplay: Assignment): MessageEmbed {
        val embed = EmbedBuilder()
        embed.title = toDisplay.name
        embed.url = toDisplay.url
        embed.field("Course ID", toDisplay.course.courseId, true)
        embed.field("Status", toDisplay.status, true)
        embed.field("Due Date", "<t:${toDisplay.due}:F>", false)
        embed.footer("ID: ${toDisplay.id}")

        return embed.build()
    }

    override fun component(toDisplay: Assignment): MessageTopLevelComponent {
        return ActionRow {
            Button.primary("edit_assignment", "📝 Edit")
            Button.secondary("delete_assignment", "🗑️ Delete")
        }
    }
}
