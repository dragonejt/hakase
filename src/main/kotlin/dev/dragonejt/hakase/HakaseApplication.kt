package dev.dragonejt.hakase

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication class HakaseApplication

fun main(args: Array<String>) {
  runApplication<HakaseApplication>(*args)
}
