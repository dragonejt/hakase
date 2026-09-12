package dev.dragonejt.hakase.osdk

import com.palantir.osdk.api.actions.ValidationResult
import dev.dragonejt.hakase.datamodels.Assignment
import dev.dragonejt.hakase_sdk.Ontology5a5029d53a2342fdA36a2cf9e72e73b9 as Ontology
import dev.dragonejt.hakase_sdk.actions.CreateAssignmentActionRequest
import dev.dragonejt.hakase_sdk.actions.CreateAssignmentActionResponse
import dev.dragonejt.hakase_sdk.actions.DeleteAssignmentActionRequest
import dev.dragonejt.hakase_sdk.actions.DeleteAssignmentActionResponse
import dev.dragonejt.hakase_sdk.actions.EditAssignmentActionRequest
import dev.dragonejt.hakase_sdk.actions.EditAssignmentActionResponse
import java.util.Optional
import org.springframework.stereotype.Repository

@Repository
class AssignmentRepository(private val ontology: Ontology) {
    fun findById(id: String): Optional<Assignment> {
        val osdkAssignment = ontology.objects().Assignment().fetch(id)
        return osdkAssignment.map { assignment ->
            Assignment(
                assignment.id().orElse(id),
                assignment.courseId().orElse(""),
                assignment.due().get(),
                assignment.name().orElse("Untitled"),
                assignment.status().orElse("Unknown"),
                assignment.url().orElse(null),
            )
        }
    }

    fun findAllByCourseId(courseId: String): List<Assignment> {
        val osdkAssignments =
            ontology
                .objects()
                .Assignment()
                .where(
                    dev.dragonejt.hakase_sdk.objectsetqueries.AssignmentFilter.courseId()
                        .eq(courseId)
                )
                .orderBy(dev.dragonejt.hakase_sdk.objectsetqueries.AssignmentOrdering.DUE_ASC)
                .fetchStream()

        return osdkAssignments
            .map { assignment ->
                Assignment(
                    assignment.id().get(),
                    assignment.courseId().orElse(""),
                    assignment.due().get(),
                    assignment.name().orElse("Untitled"),
                    assignment.status().orElse("Unknown"),
                    assignment.url().orElse(null),
                )
            }
            .toList()
    }

    fun save(entity: Assignment): Assignment {
        if (findById(entity.id).isPresent) {
            val response: EditAssignmentActionResponse =
                ontology
                    .actions()
                    .editAssignment()
                    .apply(
                        EditAssignmentActionRequest.builder()
                            .assignment(entity.id)
                            .due(entity.dueDate)
                            .name(entity.name)
                            .status(entity.status)
                            .url(entity.url ?: "")
                            .build()
                    )
            if (response.validationResult.validation.result == ValidationResult.VALID) {
                return entity
            }
        } else {
            val response: CreateAssignmentActionResponse =
                ontology
                    .actions()
                    .createAssignment()
                    .applyReturningEdits(
                        CreateAssignmentActionRequest.builder()
                            .courseId(entity.courseID)
                            .due(entity.dueDate)
                            .name(entity.name)
                            .status(entity.status)
                            .url(entity.url ?: "")
                            .build()
                    )
            if (response.validationResult.validation.result == ValidationResult.VALID) {
                val editsResult = response.actionEdits.orElse(null)
                val objectEditsOpt = editsResult?.objectEdits

                val generatedId =
                    if (objectEditsOpt != null && objectEditsOpt.isPresent) {
                        objectEditsOpt.get().addedObjects.firstOrNull()?.primaryKey?.toString()
                            ?: entity.id
                    } else {
                        entity.id
                    }

                return entity.copy(id = generatedId)
            }
        }

        throw IllegalArgumentException("Failed to save assignment")
    }

    fun deleteById(id: String) {
        val response: DeleteAssignmentActionResponse =
            ontology
                .actions()
                .deleteAssignment()
                .apply(DeleteAssignmentActionRequest.builder().assignment(id).build())

        if (response.validationResult.validation.result == ValidationResult.VALID) {
            if (response.actionEdits.isPresent) {
                return
            }
        }

        throw IllegalArgumentException("Failed to delete assignment")
    }
}
