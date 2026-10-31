package dev.dragonejt.hakase.data

import java.util.Optional
import org.springframework.data.jpa.repository.JpaRepository

interface CourseRepository : JpaRepository<Course, String> {
    override fun findById(id: String): Optional<Course>
}
