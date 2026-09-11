package dev.dragonejt.hakase.datamodels

data class Assignment(
    val id: String,
    val courseID: String,
    val dueDate: java.time.OffsetDateTime,
    val name: String,
    val status: String,
    val url: String,
)
