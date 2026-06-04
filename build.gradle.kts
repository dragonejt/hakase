plugins {
  kotlin("jvm") version "2.4.0"
  kotlin("plugin.spring") version "2.4.0"
  id("org.springframework.boot") version "4.0.6"
  id("io.spring.dependency-management") version "1.1.7"
  id("com.ncorti.ktfmt.gradle") version "0.26.0"
}

group = "dev.dragonejt"

version = "0.0.1-SNAPSHOT"

java { toolchain { languageVersion = JavaLanguageVersion.of(25) } }

repositories { mavenCentral() }

dependencies {
  implementation("org.springframework.boot:spring-boot-starter")
  implementation("org.jetbrains.kotlin:kotlin-reflect")
  developmentOnly("org.springframework.boot:spring-boot-devtools")
  annotationProcessor("org.springframework.boot:spring-boot-configuration-processor")
  testImplementation("org.springframework.boot:spring-boot-starter-test")
  testImplementation("org.jetbrains.kotlin:kotlin-test-junit5")
  testRuntimeOnly("org.junit.platform:junit-platform-launcher")

  implementation("dev.kord:kord-core:0.18.1")
  implementation("io.github.oshai:kotlin-logging-jvm:5.1.0")
}

kotlin { compilerOptions { freeCompilerArgs.addAll("-Xjsr305=strict") } }

tasks.withType<Test> { useJUnitPlatform() }
