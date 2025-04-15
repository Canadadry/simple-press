#!/usr/bin/env bash
.PHONY: help build

default: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

sqlc: ## run sqlc
	sqlc generate

sqlc-auto-reload: ## run sqlc if a .sql file change to keep go code up to date
	reflex -r '^model.+\.sql$$' sqlc generate

init: ## init database
	rm blog.sqlite
	go run main.go -action=migration

run: ## run admin server
	go test ./...
	go run main.go -action=admin

unit-test: ## unit test
	go test ./...
