package dev.dragonejt.hakase.telemetry

import io.github.oshai.kotlinlogging.KotlinLogging

abstract class LogBase {
    protected val log = KotlinLogging.logger(this.javaClass.name)
}
