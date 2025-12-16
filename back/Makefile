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

reflex: ## install reflex
	go install github.com/cespare/reflex@latest

dev: ## watch go file and restart on change
	ulimit -n 4096 ; reflex -r '\.go$$' -s -- sh -c "go run main.go -action admin"

file_add_zip: ## regenerate zip test file
	 go run tools/multipart/main.go -i tests/admin/payloads/Archive.zip -o tests/admin/payloads/file_add_zip

test: ## unit test then end2end
	go test ./...
	ulimit -n 1000 ; cd tests; madelyne test -in conf.yml

fixture: ## build fixture
	go run main.go -action=fixture
