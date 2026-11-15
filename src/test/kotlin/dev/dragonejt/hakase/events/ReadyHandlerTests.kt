package dev.dragonejt.hakase.events

import dev.dragonejt.hakase.data.CourseRepository
import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.api.trace.Tracer
import net.dv8tion.jda.api.entities.Activity
import net.dv8tion.jda.api.events.session.ReadyEvent
import net.dv8tion.jda.api.managers.Presence
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import org.mockito.Answers
import org.mockito.Mock
import org.mockito.junit.jupiter.MockitoExtension
import org.mockito.kotlin.verify
import org.mockito.kotlin.whenever

@ExtendWith(MockitoExtension::class)
class ReadyHandlerTests {

    private val tracer: Tracer = OpenTelemetry.noop().getTracer(this.javaClass.simpleName)

    private val courseId = "1063590532711972945"

    private val channelId = "1063590534528118898"

    private val courseName = "hakase"

    @Mock private lateinit var courses: CourseRepository

    @Mock(answer = Answers.RETURNS_DEEP_STUBS) private lateinit var event: ReadyEvent

    @Mock private lateinit var presence: Presence

    private lateinit var underTest: ReadyHandler

    @BeforeEach
    fun setUp() {
        whenever { event.jda.presence }.thenReturn(presence)
        whenever { event.jda.guilds.size }.thenReturn(1)
        whenever { event.jda.selfUser.name }.thenReturn(courseName)

        underTest = ReadyHandler(tracer)
    }

    @Test
    suspend fun `onEvent() updates presence status`() {
        underTest.onEvent(event)

        verify(presence).activity =
            Activity.of(
                Activity.ActivityType.CUSTOM_STATUS,
                "hakase.dragonejt.dev | in 1 courses",
            )
    }
}
