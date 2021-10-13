BINARY=./fuze
.DEFAULT_GOAL := run
test: 
	go test -v -cover -covermode=atomic ./...

run:
	go build -o ${BINARY} *.go 
	${BINARY};

add: 
	git add . 

commit:
	git commit

push:
	git push origin dev

pull:
	git pull

 