package dev.dragonejt.hakase.interactions

import net.dv8tion.jda.api.interactions.commands.build.SlashCommandData

interface SlashCommand {
    fun command(): SlashCommandData
}
