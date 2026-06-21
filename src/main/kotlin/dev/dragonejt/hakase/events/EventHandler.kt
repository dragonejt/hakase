package dev.dragonejt.hakase.events

import dev.kord.core.Kord
import dev.kord.core.event.gateway.GatewayEvent

interface EventHandler<T : GatewayEvent> {
    suspend fun register(bot: Kord)

    suspend fun handleEvent(event: T)
}
