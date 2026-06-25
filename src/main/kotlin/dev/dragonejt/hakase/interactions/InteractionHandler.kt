package dev.dragonejt.hakase.interactions

import dev.kord.core.Kord
import dev.kord.core.entity.interaction.Interaction
import dev.kord.core.event.interaction.InteractionCreateEvent

interface InteractionHandler<T : Interaction> {
    suspend fun register(bot: Kord)

    suspend fun handleInteraction(interaction: T)
}
