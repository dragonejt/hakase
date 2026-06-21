plugins {
    kotlin("jvm") version "2.4.0"
    kotlin("plugin.spring") version "2.4.0"
    id("org.springframework.boot") version "4.1.0"
    id("io.spring.dependency-management") version "1.1.7"
    id("com.diffplug.spotless") version "8.6.0"
    jacoco
    id("dev.detekt") version "2.0.0-alpha.4"
    id("io.sentry.jvm.gradle") version "6.11.0"
}

group = "dev.dragonejt"

version = "0.0.1-SNAPSHOT"

java { toolchain { languageVersion = JavaLanguageVersion.of(25) } }

repositories { mavenCentral() }

dependencies {
    implementation("org.springframework.boot:spring-boot-starter")
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("dev.kord:kord-core:0.18.1")
    implementation("io.github.oshai:kotlin-logging-jvm:5.1.0")
    implementation("io.sentry:sentry-spring-boot-4:8.44.0")
    implementation("io.sentry:sentry-async-profiler:8.44.0")
    implementation("io.sentry:sentry-opentelemetry-otlp-spring:8.44.0")
    implementation(platform("io.opentelemetry.instrumentation:opentelemetry-instrumentation-bom:2.28.1"))

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5")
    testImplementation("org.mockito.kotlin:mockito-kotlin:6.3.0")
    testImplementation("org.assertj:assertj-core:3.27.7")

    developmentOnly("org.springframework.boot:spring-boot-devtools")
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