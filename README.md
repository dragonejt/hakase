# hakase-discord
[![codacy](https://app.codacy.com/project/badge/Grade/dc0e8d6ee88549f3b6f58b5c66d39040)](https://app.codacy.com/gh/dragonejt/hakase/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![coverage](https://app.codacy.com/project/badge/Coverage/dc0e8d6ee88549f3b6f58b5c66d39040)](https://app.codacy.com/gh/dragonejt/hakase/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)
[![integration](https://github.com/dragonejt/hakase/actions/workflows/integrate.yml/badge.svg)](https://github.com/dragonejt/hakase/actions/workflows/integrate.yml)
[![delivery](https://github.com/dragonejt/hakase/actions/workflows/deliver.yml/badge.svg)](https://github.com/dragonejt/hakase/actions/workflows/deliver.yml)

hakase is a collection of helpful utilities for class chatrooms, including an assignment due date reminder, study session scheduler, and more. It is currently under development. This repository holds the Discord Bot, built with Spring Boot on Kotlin.

## Local Development
### Building and Running
To start off in VS Code, install the [Kotlin by JetBrains](https://marketplace.visualstudio.com/items?itemName=JetBrains.kotlin-server) extension. Local development with hakase-discord is relatively simple. The only command you have to run is:
```sh
./gradlew bootRun
```
You do have to have some environment variables in place. hakase does not directly read from a .env file, but you can configure environment variables or reference a .env file through IDE launch options. Otherwise, you can set environment variables locally.
```sh
ENV="development"
DISCORD_BOT_TOKEN="from Discord Dev Portal"
BACKEND_URL="from Backend API"
BACKEND_AUTH_TOKEN="from Backend API"
```
To only build the project without running the Spring Boot application, run the following:
```bash
./gradlew build
```

### Testing
For testing, the following command should be run, with the above environment variables in place:
```bash
./gradlew test
```
Tests are JUnit based, with Mockito for the mocking library. The integrate.yml GitHub Actions workflow will run these tests with coverage.

### Linting and Formatting
This project uses [detekt](https://detekt.dev/) for linting and [ktfmt](https://facebook.github.io/ktfmt/) for code formatting. The gradle build task will check both linting and code formatting, but to fix linting and formatting errors, run:
```bash
./gradlew detekt # linting
./gradlew spotlessApply # formatting
```
The `integrate.yml` GitHub Actions workflow will check for linting errors and formatting mistakes as a part of the gradle build task.

## Deployment
For deployment, hakase is built into a Docker image with [Heroku Cloud Native Buildpacks](https://github.com/heroku/cnb-builder-images), and then deployed into a container via Palantir Compute Modules.

On the deployed docker container, the same environment variables should be set, with `ENV` now being `production`.

### Continuous Delivery
hakase has a continuous delivery GitHub Actions workflow, `deliver.yml`. The steps taken are summarized:

1. Build a Docker image with the [Heroku Cloud Native Buildpacks](https://github.com/heroku/cnb-builder-images)
2. Publish the docker image to Palantir Artifacts.
3. Select the latest docker image to run in the Palantir Compute Module
4. A new Sentry release is created for monitoring with the [Sentry Release GitHub Action](https://github.com/getsentry/action-release).