plugins {
    kotlin("jvm") version "2.4.0"
    kotlin("plugin.spring") version "2.4.0"
    id("org.springframework.boot") version "4.1.0"
    id("io.spring.dependency-management") version "1.1.7"
    id("com.diffplug.spotless") version "8.6.0"
    jacoco
    id("org.sonarqube") version "7.3.1.8318"
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
    testImplementation("org.mockito.kotlin:mockito-kotlin:6.3.0")
    implementation("io.sentry:sentry-spring-boot-4:8.43.2")
    implementation("io.sentry:sentry-async-profiler:8.43.2")
}

kotlin { compilerOptions { freeCompilerArgs.addAll("-Xjsr305=strict") } }

tasks.withType<Test> {
    useJUnitPlatform() 

    reports {
        junitXml.required.set(true)
    }

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
        property("sonar.organization", "dragonejt")
        property("sonar.projectKey", "dragonejt_hakase")
        property("sonar.junit.reportPaths", "build/test-results/test")
        property("sonar.exclusions", "**/build/**")
    }
}