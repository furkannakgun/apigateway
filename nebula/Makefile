PROJECT_NAME = $(notdir ${PWD})
# Convert to lowercase
PROJECT_NAME := $(shell echo ${PROJECT_NAME} | tr A-Z a-z)

ENV ?= dev
export ENV
# Log to file to cwd instead of /var/log/
NAAP_STRUCTLOG_FILE ?= naap-service.log
export NAAP_STRUCTLOG_FILE

# Enable flask debug (hot reload with file changes)
FLASK_ENV ?= development

TARGET_TEST=$(if ${TEST},${TEST},tests/)
COV=$(if ${TEST},,--cov=app --cov-report=term --cov-report=xml:coverage.xml --junitxml=report.xml)

# Nexus credentials for building docker image
# TODO: use `--secret` once the docker daemon on Jenkins supports buildkit
# DOCKER_BUILD_ARGS=--secret id=nexus_user,env=NEXUS_USER --secret=id=nexus_password,env=NEXUS_PASSWORD
DOCKER_BUILD_ARGS=--build-arg NEXUS_USER=${NEXUS_USER} --build-arg NEXUS_PASSWORD=${NEXUS_PASSWORD}

help:
	@echo "help"
	@echo "       Print this help"
	@echo "build"
	@echo "       Build docker image. \
export NEXUS_USER and NEXUS_PASSWORD to authenticate to artifact repository."
	@echo "check"
	@echo "       Run static analysis"
	@echo "test"
	@echo "       Run unit tests with coverage. \
Use TEST=/path/to/test to run a specific test"
	@echo "run"
	@echo "       Run the application. \
Development mode is enabled by default: use FLASK_ENV=production to disable it"

check:
	@go build -v && rm ./nebula

build:
	@docker build . -t ${PROJECT_NAME}:${ENV} ${DOCKER_BUILD_ARGS}

test:
	@echo "Tests"
	@echo "---------------------------------------------------------------"
	@go test -coverpkg=./... -covermode=count -coverprofile coverage ./...
	@echo "Coverage"
	@echo "---------------------------------------------------------------"
	@go tool cover -func=coverage

# Visual results to help devs
devtest:
	@go test -coverpkg=./... -covermode=count -coverprofile coverage ./...
	@go tool cover -html=coverage

format:
	@gofmt -l .

run:
	@go run .


.PHONY: help check build test run
