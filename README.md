# hakase-discord
[![godoc](https://pkg.go.dev/badge/github.com/dragonejt/hakase.svg)](https://pkg.go.dev/github.com/dragonejt/hakase)
[![integration](https://github.com/dragonejt/hakase/actions/workflows/integrate.yml/badge.svg)](https://github.com/dragonejt/hakase/actions/workflows/integrate.yml)
[![delivery](https://github.com/dragonejt/hakase/actions/workflows/deliver.yml/badge.svg)](https://github.com/dragonejt/hakase/actions/workflows/deliver.yml)
[![codecov](https://codecov.io/gh/dragonejt/hakase/graph/badge.svg?token=7MEF3IHI00)](https://codecov.io/gh/dragonejt/hakase)
[![go report card](https://goreportcard.com/badge/github.com/dragonejt/hakase)](https://goreportcard.com/report/github.com/dragonejt/hakase)

hakase is a collection of helpful utilities for class chatrooms, including an assignment due date reminder, study session scheduler, and more. It is currently under development. This repository holds the Discord Bot, built with Go.

## Local Development
### Building and Running
Local development with hakase-discord is relatively simple. The only command you have to run is:
```sh
go run hakase-discord.go
```
You do have to have some environment variables in place. hakase does not directly read from a .env file, but you can configure environment variables or reference a .env file through IDE launch options. Otherwise, you can set environment variables locally.
```sh
ENV="development"
DISCORD_BOT_TOKEN="from Discord Dev Portal"
CF_ACCOUNT_ID="from Cloudflare"
CF_API_TOKEN="from Cloudflare User Tokens"
D1_DATABASE_ID="from Cloudflare D1"
```

### Testing
For testing, the following command should be run, with the above environment variables in place:
```bash
go test ./...
```
This uses Go's built-in test runner which will discover and test all `_test.go` files. The integrate.yml GitHub Actions workflow will run these tests with code coverage (`-coverpkg=./... -coverprofile=coverage.txt`).

If you are using VS Code, the [VS Code Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) will enable automatic test discovery and running in the Testing sidebar.

### Linting and Formatting
Go and the [VS Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) automatically performs linting and formatting on save. 
The `integrate.yml` GitHub Actions workflow will check for linting errors and formatting mistakes with [golangci-lint](https://github.com/golangci/golangci-lint-action).

## Deployment
For deployment, hakase is built into a Docker image with [Heroku Cloud Native Buildpacks](https://github.com/heroku/cnb-builder-images), and then deployed into a container via Palantir Compute Modules.

On the deployed docker container, the same environment variables should be set, with `ENV` now being `production`.

### Continuous Delivery
hakase has a continuous delivery GitHub Actions workflow, `deliver.yml`. The steps taken are summarized:

1. Build a Docker image with the [Heroku Cloud Native Buildpacks](https://github.com/heroku/cnb-builder-images)
2. Publish the docker image to Palantir Artifacts.
3. Select the latest docker image to run in the Palantir Compute Module
4. A new Sentry release is created for monitoring with the [Sentry Release GitHub Action](https://github.com/getsentry/action-release).