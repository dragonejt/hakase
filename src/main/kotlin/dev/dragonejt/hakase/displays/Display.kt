package dev.dragonejt.hakase.displays

import net.dv8tion.jda.api.components.MessageTopLevelComponent
import net.dv8tion.jda.api.entities.MessageEmbed

interface Display<T> {
    fun embed(toDisplay: T): MessageEmbed

    fun component(toDisplay: T): MessageTopLevelComponent
}
