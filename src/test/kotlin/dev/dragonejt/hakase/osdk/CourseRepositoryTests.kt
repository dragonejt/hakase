package dev.dragonejt.hakase.osdk

import dev.dragonejt.hakase.datamodels.Course
import dev.dragonejt.hakase_sdk.Ontology5a5029d53a2342fdA36a2cf9e72e73b9 as Ontology
import dev.dragonejt.hakase_sdk.Ontology5a5029d53a2342fdA36a2cf9e72e73b9Actions as OntologyActions
import dev.dragonejt.hakase_sdk.Ontology5a5029d53a2342fdA36a2cf9e72e73b9BaseObjectSets as OntologyBaseObjectSets
import dev.dragonejt.hakase_sdk.objectsets.CourseObjectSet
import java.util.Optional
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertAll
import org.junit.jupiter.api.extension.ExtendWith
import org.mockito.Mock
import org.mockito.junit.jupiter.MockitoExtension
import org.mockito.kotlin.any
import org.mockito.kotlin.mock
import org.mockito.kotlin.verify
import org.mockito.kotlin.whenever

@ExtendWith(MockitoExtension::class)
class CourseRepositoryTests {

    @Mock private lateinit var ontology: Ontology
    private lateinit var underTest: CourseRepository

    @BeforeEach
    fun setUp() {
        underTest = CourseRepository(ontology)
    }

    @Test
    @DisplayName("CourseRepository saves new course")
    fun testNewCourseSave() {
        val course = Course("test_course_id", "channel_id", "group_id")

        val mockObjects = mock<OntologyBaseObjectSets>()
        whenever { ontology.objects() }.thenReturn(mockObjects)

        val mockCourse = mock<CourseObjectSet>()
        whenever { mockObjects.Course() }.thenReturn(mockCourse)
        whenever { mockCourse.fetch(course.courseId) }.thenReturn(Optional.empty())

        val mockActions = mock<OntologyActions>()
        whenever { ontology.actions() }.thenReturn(mockActions)

        val mockCreateCourse = mock<dev.dragonejt.hakase_sdk.actions.CreateCourseAction>()
        whenever { mockActions.createCourse() }.thenReturn(mockCreateCourse)

        val mockResponse =
            org.mockito.Mockito.mock(
                dev.dragonejt.hakase_sdk.actions.CreateCourseActionResponse::class.java,
                org.mockito.Mockito.RETURNS_DEEP_STUBS,
            )
        whenever {
                mockCreateCourse.apply(
                    any<dev.dragonejt.hakase_sdk.actions.CreateCourseActionRequest>()
                )
            }
            .thenReturn(mockResponse)

        whenever { mockResponse.validationResult.validation.result }
            .thenReturn(com.palantir.osdk.api.actions.ValidationResult.VALID)
        val mockEdits = mock<dev.dragonejt.hakase_sdk.actions.CreateCourseActionEditsResult>()
        whenever { mockResponse.actionEdits }.thenReturn(Optional.of(mockEdits))

        val result = underTest.save(course)

        assertAll(
            { verify(mockCourse).fetch(course.courseId) },
            {
                verify(mockCreateCourse)
                    .apply(any<dev.dragonejt.hakase_sdk.actions.CreateCourseActionRequest>())
            },
            { assertThat(result).isEqualTo(course) },
        )
    }

    @Test
    @DisplayName("CourseRepository updates existing course")
    fun testExistingCourseSave() {
        val course = Course("test_course_id", "channel_id", "group_id")

        val mockObjects = mock<OntologyBaseObjectSets>()
        whenever { ontology.objects() }.thenReturn(mockObjects)

        val mockCourse = mock<CourseObjectSet>()
        whenever { mockObjects.Course() }.thenReturn(mockCourse)

        val mockOsdkCourse = mock<dev.dragonejt.hakase_sdk.objects.Course>()
        whenever { mockOsdkCourse.courseId() }.thenReturn(Optional.of(course.courseId))
        whenever { mockOsdkCourse.notifyChannel() }.thenReturn(Optional.of("channel_id"))
        whenever { mockOsdkCourse.notifyGroup() }.thenReturn(Optional.of("group_id"))
        whenever { mockCourse.fetch(course.courseId) }.thenReturn(Optional.of(mockOsdkCourse))

        val mockActions = mock<OntologyActions>()
        whenever { ontology.actions() }.thenReturn(mockActions)

        val mockEditCourse = mock<dev.dragonejt.hakase_sdk.actions.EditCourseAction>()
        whenever { mockActions.editCourse() }.thenReturn(mockEditCourse)

        val mockResponse =
            org.mockito.Mockito.mock(
                dev.dragonejt.hakase_sdk.actions.EditCourseActionResponse::class.java,
                org.mockito.Mockito.RETURNS_DEEP_STUBS,
            )
        whenever {
                mockEditCourse.apply(
                    any<dev.dragonejt.hakase_sdk.actions.EditCourseActionRequest>()
                )
            }
            .thenReturn(mockResponse)

        whenever { mockResponse.validationResult.validation.result }
            .thenReturn(com.palantir.osdk.api.actions.ValidationResult.VALID)
        val mockEdits = mock<dev.dragonejt.hakase_sdk.actions.EditCourseActionEditsResult>()
        whenever { mockResponse.actionEdits }.thenReturn(Optional.of(mockEdits))

        val result = underTest.save(course)

        assertAll(
            { verify(mockCourse).fetch(course.courseId) },
            {
                verify(mockEditCourse)
                    .apply(any<dev.dragonejt.hakase_sdk.actions.EditCourseActionRequest>())
            },
            { assertThat(result).isEqualTo(course) },
        )
    }

    @Test
    @DisplayName("CourseRepository finds course by ID")
    fun testFindById() {
        val courseId = "test_course_id"

        val mockObjects = mock<OntologyBaseObjectSets>()
        whenever { ontology.objects() }.thenReturn(mockObjects)

        val mockCourse = mock<CourseObjectSet>()
        whenever { mockObjects.Course() }.thenReturn(mockCourse)

        val mockOsdkCourse = mock<dev.dragonejt.hakase_sdk.objects.Course>()
        whenever { mockOsdkCourse.courseId() }.thenReturn(Optional.of(courseId))
        whenever { mockOsdkCourse.notifyChannel() }.thenReturn(Optional.of("channel_id"))
        whenever { mockOsdkCourse.notifyGroup() }.thenReturn(Optional.of("group_id"))

        whenever { mockCourse.fetch(courseId) }.thenReturn(Optional.of(mockOsdkCourse))

        val result = underTest.findById(courseId)

        assertAll(
            { verify(mockCourse).fetch(courseId) },
            { assertThat(result).isPresent() },
            { assertThat(result.get().courseId).isEqualTo(courseId) },
            { assertThat(result.get().notifyChannel).isEqualTo("channel_id") },
            { assertThat(result.get().notifyGroup).isEqualTo("group_id") },
        )
    }

    @Test
    @DisplayName("CourseRepository deletes course by ID")
    fun testDeleteById() {
        val courseId = "test_course_id"

        val mockActions = mock<OntologyActions>()
        whenever { ontology.actions() }.thenReturn(mockActions)

        val mockDeleteCourse = mock<dev.dragonejt.hakase_sdk.actions.DeleteCourseAction>()
        whenever { mockActions.deleteCourse() }.thenReturn(mockDeleteCourse)

        val mockResponse =
            org.mockito.Mockito.mock(
                dev.dragonejt.hakase_sdk.actions.DeleteCourseActionResponse::class.java,
                org.mockito.Mockito.RETURNS_DEEP_STUBS,
            )
        whenever {
                mockDeleteCourse.apply(
                    any<dev.dragonejt.hakase_sdk.actions.DeleteCourseActionRequest>()
                )
            }
            .thenReturn(mockResponse)

        whenever { mockResponse.validationResult.validation.result }
            .thenReturn(com.palantir.osdk.api.actions.ValidationResult.VALID)
        val mockEdits = mock<dev.dragonejt.hakase_sdk.actions.DeleteCourseActionEditsResult>()
        whenever { mockResponse.actionEdits }.thenReturn(Optional.of(mockEdits))

        underTest.deleteById(courseId)

        verify(mockDeleteCourse)
            .apply(any<dev.dragonejt.hakase_sdk.actions.DeleteCourseActionRequest>())
    }
}
