run: build
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
	@templ generate --watch --proxy=http://localhost:3000

build:	
	@templ generate views
	@go build -o bin/app main.go