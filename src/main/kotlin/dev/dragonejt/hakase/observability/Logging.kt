package dev.dragonejt.hakase.observability

import io.github.oshai.kotlinlogging.KLogger
import io.github.oshai.kotlinlogging.KotlinLogging

interface Logging {
    val logger: KLogger
        get() {
            val kClass = this::class
            val owner =
                if (kClass.isCompanion) kClass.java.enclosingClass ?: kClass.java else kClass.java
            return KotlinLogging.logger(owner.name)
        }
}
