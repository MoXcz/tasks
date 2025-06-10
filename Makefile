run: build
	@./bin/tasks

build:
	@go build -o bin/tasks main.go

add: build
	@./bin/tasks add "example description"

list: build
	@./bin/tasks list
