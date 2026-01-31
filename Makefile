.PHONY: help run build test lint clean

APP_NAME=gateway
BUILD_DIR=bin
CONFIG_DIR=configs
HOME_DIR := $(shell echo ~)

export GOCACHE=$(HOME_DIR)/.cache/go-build
export GOTMPDIR=/tmp
export GOPATH=$(HOME_DIR)/go
export GOMODCACHE=$(GOPATH)/pkg/mod
export GOLANGCI_LINT_CACHE=$(HOME_DIR)/.cache/golangci-lint

help: ## Показать эту справку
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Запустить приложение
	go run cmd/$(APP_NAME)/main.go

build: ## Собрать приложение
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/$(APP_NAME)/main.go

test: ## Запустить тесты
	go test -v -race -coverprofile=coverage.out ./...

lint: ## Запустить линтер
	mkdir -p $(GOLANGCI_LINT_CACHE)
	@rm -f $(GOLANGCI_LINT_CACHE)/golangci-lint.lock
	golangci-lint run

clean: ## Удалить артефакты сборки
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

deps: ## Установить зависимости
	go mod download
	go mod tidy

fmt: ## Форматировать код
	go fmt ./...

.DEFAULT_GOAL := help