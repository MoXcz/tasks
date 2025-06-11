run: build
	@./bin/tasks add "example description" && \
		./bin/tasks list && \
		./bin/tasks add "another example" && \
		./bin/tasks complete 1 && \
		./bin/tasks list -a

build:
	@go build -o bin/tasks main.go

add: build
	@./bin/tasks add "example description"

list: build
	@./bin/tasks list

clean:
	rm -r ./bin

