package dev.dragonejt.hakase.telemetry

import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.api.trace.Tracer
import io.sentry.opentelemetry.otlp.OpenTelemetryOtlpEventProcessor
import org.springframework.beans.factory.annotation.Value
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@Configuration
class TelemetryConfig {

    @Bean
    fun eventProcessor(): OpenTelemetryOtlpEventProcessor {
        return OpenTelemetryOtlpEventProcessor()
    }

    @Bean
    fun tracer(
        openTelemetry: OpenTelemetry,
        @Value($$"${spring.application.name}") appName: String,
    ): Tracer {
        return openTelemetry.tracerBuilder(appName).build()
    }
}
