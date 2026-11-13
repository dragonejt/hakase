package dev.dragonejt.hakase.data

import jakarta.persistence.Entity
import jakarta.persistence.GeneratedValue
import jakarta.persistence.GenerationType
import jakarta.persistence.Id
import jakarta.persistence.JoinColumn
import jakarta.persistence.ManyToOne
import java.time.Instant

@Entity
data class Assignment(
    @Id @GeneratedValue(strategy = GenerationType.IDENTITY) val id: Long?,
    @ManyToOne @JoinColumn(name = "course_id", nullable = false) val course: Course,
    val due: Instant,
    val name: String,
    val status: String,
    val url: String? = null,
)
