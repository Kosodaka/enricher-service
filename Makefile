
BINARY_NAME=enricher-service
.DEFAULT_GOAL:=run
#to make build use git bash or linux console
#build app for windows and linux os and save in ./target directory
build:
	 GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-windows.exe cmd/app/main.go
	 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux.exe cmd/app/main.go
 #run windows build
 run-windows: build
	${BINARY_NAME}-windows
#run linux build
run-linux: build
	${BINARY_NAME}-linux
test:
	go test -cover ./...

docker-up:
	docker-compose up -d
migrate:
	goose -dir ./migrations postgres "postgres://admin:qwerty@localhost:5432/human?sslmode=disable" up
migrate-down:
	goose -dir ./migrations postgres "postgres://admin:qwerty@localhost:5432/human?sslmode=disable" down
# clean builds in ./target
clean:
	go clean
	rm ${BINARY_NAME}-windows
	rm ${BINARY_NAME}-linux
# install dependencies
dep:
	go mod download