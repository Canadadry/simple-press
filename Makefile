#!/usr/bin/env bash
.PHONY: help build

default: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

sqlc: ## run sqlc
	sqlc generate

sqlc-auto-reload: ## run sqlc if a .sql file change to keep go code up to date
	reflex -r '^model.+\.sql$$' sqlc generate

blog.sqlite: ## init database
	go run main.go -action=migration

demo: ## run admin server with demo blog
	go run main.go -action=admin -db-url=example/starter.blog

run: blog.sqlite ## run admin server
	go test ./...
	go run main.go -action=admin

runp: blog.sqlite ## run public server
	go test ./...
	go run main.go

unit-test: ## unit test
	go test ./...
