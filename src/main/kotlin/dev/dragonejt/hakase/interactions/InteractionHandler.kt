package dev.dragonejt.hakase.interactions

import dev.kord.core.Kord
import dev.kord.core.entity.interaction.Interaction

interface InteractionHandler<T : Interaction> {
    suspend fun register(bot: Kord)

    suspend fun handleInteraction(interaction: T)
}
