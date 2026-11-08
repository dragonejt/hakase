plugins {
    kotlin("jvm") version "2.4.0"
    kotlin("plugin.spring") version "2.4.0"
    kotlin("plugin.jpa") version "2.4.0"
    id("org.springframework.boot") version "4.1.0"
    id("io.spring.dependency-management") version "1.1.7"
    id("com.diffplug.spotless") version "8.10.3"
    jacoco
    id("dev.detekt") version "2.0.0-alpha.4"
    id("io.sentry.jvm.gradle") version "6.11.0"
    id("org.sonarqube") version "7.5.0.8588"
}

group = "dev.dragonejt"

version = "0.0.1-SNAPSHOT"

extra["opentelemetry.version"] = "1.65.0"

java { toolchain { languageVersion = JavaLanguageVersion.of(25) } }

repositories {
    mavenCentral()     
}

dependencies {
    // Spring Boot
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    implementation("org.jetbrains.kotlin:kotlin-reflect")

    implementation("net.dv8tion:JDA:6.7.0")
    implementation("club.minnced:jda-ktx:0.15.0")
    runtimeOnly("org.postgresql:postgresql")
    implementation("io.github.natty-parser:natty:1.1.4")

    // OpenTelemetry
    implementation("io.github.oshai:kotlin-logging-jvm:8.0.4")
    implementation("io.sentry:sentry-spring-boot-4:8.58.0")
    implementation("io.sentry:sentry-async-profiler:8.58.0")
    implementation("io.sentry:sentry-opentelemetry-otlp-spring:8.58.0")
    implementation("io.opentelemetry.instrumentation:opentelemetry-spring-boot-starter")
    implementation(platform("io.opentelemetry.instrumentation:opentelemetry-instrumentation-bom:2.31.1"))
    

    testImplementation("org.springframework.boot:spring-boot-starter-data-jpa-test")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5")
    testImplementation("org.mockito.kotlin:mockito-kotlin:6.4.0")
    testImplementation("org.assertj:assertj-core:3.27.7")

    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
    annotationProcessor("org.springframework.boot:spring-boot-configuration-processor")

}

kotlin { compilerOptions { freeCompilerArgs.addAll("-Xjsr305=strict") } }

tasks.withType<Test> {
    useJUnitPlatform() 
    finalizedBy(tasks.jacocoTestReport)
}

spotless {
    kotlin {
        ktfmt("0.63").kotlinlangStyle()
    }
}

tasks.jacocoTestReport {
    dependsOn(tasks.test)

    reports {
        xml.required.set(true)
    }
}

sonar {
  properties {
    property("sonar.projectKey", "dragonejt_hakase")
    property("sonar.organization", "dragonejt")
  }
}