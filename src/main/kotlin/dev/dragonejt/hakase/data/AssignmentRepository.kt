package dev.dragonejt.hakase.data

import java.util.Optional
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface AssignmentRepository : JpaRepository<Assignment, Long> {
    override fun findById(id: Long): Optional<Assignment>
}
