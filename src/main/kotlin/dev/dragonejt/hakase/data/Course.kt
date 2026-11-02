package dev.dragonejt.hakase.data

import jakarta.persistence.Entity
import jakarta.persistence.FetchType
import jakarta.persistence.Id
import jakarta.persistence.OneToMany

@Entity
data class Course(
    @Id val courseId: String,
    val notifyChannel: String,
    val notifyGroup: String,
    @OneToMany(mappedBy = "course", fetch = FetchType.EAGER)
    val assignments: MutableList<Assignment> = mutableListOf(),
)
