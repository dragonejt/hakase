package dev.dragonejt.hakase.interactions

import net.dv8tion.jda.api.interactions.commands.build.CommandData

interface ApplicationCommand {
    fun command(): CommandData
}
