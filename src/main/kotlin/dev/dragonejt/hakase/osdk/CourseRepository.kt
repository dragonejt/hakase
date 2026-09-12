package dev.dragonejt.hakase.osdk

import com.palantir.osdk.api.actions.ValidationResult
import dev.dragonejt.hakase.datamodels.Course
import dev.dragonejt.hakase_sdk.Ontology5a5029d53a2342fdA36a2cf9e72e73b9 as Ontology
import dev.dragonejt.hakase_sdk.actions.CreateCourseActionRequest
import dev.dragonejt.hakase_sdk.actions.CreateCourseActionResponse
import dev.dragonejt.hakase_sdk.actions.DeleteCourseActionRequest
import dev.dragonejt.hakase_sdk.actions.DeleteCourseActionResponse
import dev.dragonejt.hakase_sdk.actions.EditCourseActionRequest
import dev.dragonejt.hakase_sdk.actions.EditCourseActionResponse
import java.util.Optional
import org.springframework.data.repository.CrudRepository
import org.springframework.stereotype.Repository

@Repository
@Suppress("TooManyFunctions")
class CourseRepository(private val ontology: Ontology) : CrudRepository<Course, String> {
    override fun <S : Course> save(entity: S): S {
        if (findById(entity.courseId).isPresent) {
            val response: EditCourseActionResponse =
                ontology
                    .actions()
                    .editCourse()
                    .apply(
                        EditCourseActionRequest.builder()
                            .course(entity.courseId)
                            .notifyChannel(entity.notifyChannel)
                            .notifyGroup(entity.notifyGroup)
                            .build()
                    )
            if (response.validationResult.validation.result == ValidationResult.VALID) {
                if (response.actionEdits.isPresent) {
                    return entity
                }
            }
        } else {
            val response: CreateCourseActionResponse =
                ontology
                    .actions()
                    .createCourse()
                    .apply(
                        CreateCourseActionRequest.builder()
                            .notifyChannel(entity.notifyChannel)
                            .notifyGroup(entity.notifyGroup)
                            .courseId(entity.courseId)
                            .build()
                    )
            if (response.validationResult.validation.result == ValidationResult.VALID) {
                if (response.actionEdits.isPresent) {
                    return entity
                }
            }
        }

        throw IllegalArgumentException("Failed to save course")
    }

    override fun <S : Course> saveAll(entities: Iterable<S>): Iterable<S> {
        throw UnsupportedOperationException("Not implemented")
    }

    override fun findById(id: String): Optional<Course> {
        val osdkCourse = ontology.objects().Course().fetch(id)
        return osdkCourse.map { course ->
            Course(
                course.courseId().get(),
                course.notifyChannel().get(),
                course.notifyGroup().get(),
            )
        }
    }

    override fun existsById(id: String): Boolean {
        throw UnsupportedOperationException("Not implemented")
    }

    override fun findAll(): Iterable<Course> {
        throw UnsupportedOperationException("Not implemented")
    }

    override fun findAllById(ids: Iterable<String>): Iterable<Course> {
        throw UnsupportedOperationException("Not implemented")
    }

    override fun count(): Long {
        throw UnsupportedOperationException("Not implemented")
    }

    override fun deleteById(id: String): Unit {
        val response: DeleteCourseActionResponse =
            ontology
                .actions()
                .deleteCourse()
                .apply(DeleteCourseActionRequest.builder().course(id).build())

        if (response.validationResult.validation.result == ValidationResult.VALID) {
            if (response.actionEdits.isPresent) {
                return
            }
        }

        throw IllegalArgumentException("Failed to delete course")
    }

    override fun delete(entity: Course): Unit {
        deleteById(entity.courseId)
    }

    override fun deleteAllById(ids: Iterable<String>): Unit {
        throw UnsupportedOperationException("Not implemented")
    }

    override fun deleteAll(entities: Iterable<Course>): Unit {
        throw UnsupportedOperationException("Not implemented")
    }

    override fun deleteAll(): Unit {
        throw UnsupportedOperationException("Not implemented")
    }
}
