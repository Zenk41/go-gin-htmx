run: tailwind build-in-nix
	@./bin/app

install: 
	@go install github.com/a-h/templ/@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D

tailwind:
	@tailwindcss -i styles/input.css -o public/globals.css

templ:
	@templ generate --watch --proxy=http://localhost:8080

build:	
	@templ generate views
	@go build -o bin/app main.go

nix-templ:
	@nix run github.com/a-h/templ generate --watch --proxy=http://localhost:8080

build-in-nix:
	# for anything other than nix	
	# @templ generate views
	
	# for nix 
	@nix run github:a-h/templ generate views
	@go build -o bin/app main.go