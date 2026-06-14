package dev.dragonejt.hakase.events

import dev.kord.core.Kord
import dev.kord.core.event.gateway.GatewayEvent

interface EventHandler<T : GatewayEvent> {
    fun register(bot: Kord)

    fun handleEvent(event: T)
}
