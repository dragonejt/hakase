package dev.dragonejt.hakase.osdk

import dev.dragonejt.hakase.datamodels.Assignment
import dev.dragonejt.hakase_sdk.Ontology5a5029d53a2342fdA36a2cf9e72e73b9 as Ontology
import dev.dragonejt.hakase_sdk.Ontology5a5029d53a2342fdA36a2cf9e72e73b9Actions as OntologyActions
import dev.dragonejt.hakase_sdk.Ontology5a5029d53a2342fdA36a2cf9e72e73b9BaseObjectSets as OntologyBaseObjectSets
import dev.dragonejt.hakase_sdk.objectsets.AssignmentObjectSet
import java.time.OffsetDateTime
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
class AssignmentRepositoryTests {
    @Mock private lateinit var ontology: Ontology
    private lateinit var underTest: AssignmentRepository

    @BeforeEach
    fun setUp() {
        underTest = AssignmentRepository(ontology)
    }

    @Test
    @DisplayName("AssignmentRepository saves new assignment")
    fun testNewAssignmentSave() {
        val assignment =
            Assignment(
                "test_assignment_id",
                "test_course_id",
                OffsetDateTime.now(),
                "test_assignment_name",
                "test_assignment_status",
                "test_assignment_url",
            )
        val mockObjects = mock<OntologyBaseObjectSets>()
        whenever { ontology.objects() }.thenReturn(mockObjects)

        val mockAssignment = mock<AssignmentObjectSet>()
        whenever { mockObjects.Assignment() }.thenReturn(mockAssignment)
        whenever { mockAssignment.fetch(assignment.id) }.thenReturn(Optional.empty())

        val mockActions = mock<OntologyActions>()
        whenever { ontology.actions() }.thenReturn(mockActions)

        val mockCreateAssignment = mock<dev.dragonejt.hakase_sdk.actions.CreateAssignmentAction>()
        whenever { mockActions.createAssignment() }.thenReturn(mockCreateAssignment)

        val mockResponse =
            org.mockito.Mockito.mock(
                dev.dragonejt.hakase_sdk.actions.CreateAssignmentActionResponse::class.java,
                org.mockito.Mockito.RETURNS_DEEP_STUBS,
            )
        whenever {
                mockCreateAssignment.applyReturningEdits(
                    any<dev.dragonejt.hakase_sdk.actions.CreateAssignmentActionRequest>()
                )
            }
            .thenReturn(mockResponse)
        whenever { mockResponse.validationResult.validation.result }
            .thenReturn(com.palantir.osdk.api.actions.ValidationResult.VALID)
        val mockEdits = mock<dev.dragonejt.hakase_sdk.actions.CreateAssignmentActionEditsResult>()
        whenever { mockEdits.objectEdits }.thenReturn(Optional.empty())
        whenever { mockResponse.actionEdits }.thenReturn(Optional.of(mockEdits))

        val result = underTest.save(assignment)

        assertAll(
            { verify(mockAssignment).fetch(assignment.id) },
            {
                verify(mockCreateAssignment)
                    .applyReturningEdits(
                        any<dev.dragonejt.hakase_sdk.actions.CreateAssignmentActionRequest>()
                    )
            },
            { assertThat(result).isEqualTo(assignment) },
        )
    }

    @Test
    @DisplayName("AssignmentRepository updates existing assignment")
    fun testUpdateExistingAssignment() {
        val assignment =
            Assignment(
                "test_assignment_id",
                "test_course_id",
                OffsetDateTime.now(),
                "test_assignment_name",
                "test_assignment_status",
                "test_assignment_url",
            )
        val mockObjects = mock<OntologyBaseObjectSets>()
        whenever { ontology.objects() }.thenReturn(mockObjects)

        val mockAssignment = mock<AssignmentObjectSet>()
        whenever { mockObjects.Assignment() }.thenReturn(mockAssignment)

        val mockOsdkAssignment = mock<dev.dragonejt.hakase_sdk.objects.Assignment>()
        whenever { mockOsdkAssignment.id() }.thenReturn(Optional.of(assignment.id))
        whenever { mockOsdkAssignment.courseId() }.thenReturn(Optional.of(assignment.courseID))
        whenever { mockOsdkAssignment.due() }.thenReturn(Optional.of(assignment.dueDate))
        whenever { mockOsdkAssignment.name() }.thenReturn(Optional.of(assignment.name))
        whenever { mockOsdkAssignment.status() }.thenReturn(Optional.of(assignment.status))
        whenever { mockOsdkAssignment.url() }.thenReturn(Optional.ofNullable(assignment.url))
        whenever { mockAssignment.fetch(assignment.id) }.thenReturn(Optional.of(mockOsdkAssignment))

        val mockActions = mock<OntologyActions>()
        whenever { ontology.actions() }.thenReturn(mockActions)

        val mockEditAssignment = mock<dev.dragonejt.hakase_sdk.actions.EditAssignmentAction>()
        whenever { mockActions.editAssignment() }.thenReturn(mockEditAssignment)

        val mockResponse =
            org.mockito.Mockito.mock(
                dev.dragonejt.hakase_sdk.actions.EditAssignmentActionResponse::class.java,
                org.mockito.Mockito.RETURNS_DEEP_STUBS,
            )
        whenever {
                mockEditAssignment.apply(
                    any<dev.dragonejt.hakase_sdk.actions.EditAssignmentActionRequest>()
                )
            }
            .thenReturn(mockResponse)
        whenever { mockResponse.validationResult.validation.result }
            .thenReturn(com.palantir.osdk.api.actions.ValidationResult.VALID)

        val result = underTest.save(assignment)

        assertAll(
            { verify(mockAssignment).fetch(assignment.id) },
            {
                verify(mockEditAssignment)
                    .apply(any<dev.dragonejt.hakase_sdk.actions.EditAssignmentActionRequest>())
            },
            { assertThat(result).isEqualTo(assignment) },
        )
    }

    @Test
    @DisplayName("AssignmentRepository finds assignment by ID")
    fun testFindById() {
        val assignment =
            Assignment(
                "test_assignment_id",
                "test_course_id",
                OffsetDateTime.now(),
                "test_assignment_name",
                "test_assignment_status",
                "test_assignment_url",
            )
        val mockObjects = mock<OntologyBaseObjectSets>()
        whenever { ontology.objects() }.thenReturn(mockObjects)

        val mockAssignment = mock<AssignmentObjectSet>()
        whenever { mockObjects.Assignment() }.thenReturn(mockAssignment)

        val mockOsdkAssignment = mock<dev.dragonejt.hakase_sdk.objects.Assignment>()
        whenever { mockOsdkAssignment.id() }.thenReturn(Optional.of(assignment.id))
        whenever { mockOsdkAssignment.courseId() }.thenReturn(Optional.of(assignment.courseID))
        whenever { mockOsdkAssignment.due() }.thenReturn(Optional.of(assignment.dueDate))
        whenever { mockOsdkAssignment.name() }.thenReturn(Optional.of(assignment.name))
        whenever { mockOsdkAssignment.status() }.thenReturn(Optional.of(assignment.status))
        whenever { mockOsdkAssignment.url() }.thenReturn(Optional.ofNullable(assignment.url))

        whenever { mockAssignment.fetch(assignment.id) }.thenReturn(Optional.of(mockOsdkAssignment))

        val result = underTest.findById(assignment.id)

        assertAll(
            { verify(mockAssignment).fetch(assignment.id) },
            { assertThat(result).isPresent() },
            { assertThat(result.get().id).isEqualTo(assignment.id) },
            { assertThat(result.get().courseID).isEqualTo(assignment.courseID) },
            { assertThat(result.get().dueDate).isEqualTo(assignment.dueDate) },
            { assertThat(result.get().name).isEqualTo(assignment.name) },
            { assertThat(result.get().status).isEqualTo(assignment.status) },
            { assertThat(result.get().url).isEqualTo(assignment.url) },
        )
    }

    @Test
    @DisplayName("AssignmentRepository deletes assignment by ID")
    fun testDeleteById() {
        val assignmentId = "test_assignment_id"

        val mockActions = mock<OntologyActions>()
        whenever { ontology.actions() }.thenReturn(mockActions)

        val mockDeleteAssignment = mock<dev.dragonejt.hakase_sdk.actions.DeleteAssignmentAction>()
        whenever { mockActions.deleteAssignment() }.thenReturn(mockDeleteAssignment)

        val mockResponse =
            org.mockito.Mockito.mock(
                dev.dragonejt.hakase_sdk.actions.DeleteAssignmentActionResponse::class.java,
                org.mockito.Mockito.RETURNS_DEEP_STUBS,
            )
        whenever {
                mockDeleteAssignment.apply(
                    any<dev.dragonejt.hakase_sdk.actions.DeleteAssignmentActionRequest>()
                )
            }
            .thenReturn(mockResponse)

        whenever { mockResponse.validationResult.validation.result }
            .thenReturn(com.palantir.osdk.api.actions.ValidationResult.VALID)
        val mockEdits = mock<dev.dragonejt.hakase_sdk.actions.DeleteAssignmentActionEditsResult>()
        whenever { mockResponse.actionEdits }.thenReturn(Optional.of(mockEdits))

        underTest.deleteById(assignmentId)

        verify(mockDeleteAssignment)
            .apply(any<dev.dragonejt.hakase_sdk.actions.DeleteAssignmentActionRequest>())
    }
}
