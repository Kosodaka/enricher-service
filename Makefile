
BINARY_NAME=enricher-service
.DEFAULT_GOAL:=run
#to make build use git bash or linux console
#build app for windows and linux os and save in ./ directory
build:
	 GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-windows.exe cmd/app/main.go
	 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux.exe cmd/app/main.go
 #run windows build
 run-windows: build
	${BINARY_NAME}-windows
#run linux build
run-linux: build
	${BINARY_NAME}-linux
# install dependencies
dep:
	go mod download
docker-up:
	docker-compose up -d
test:
	go test -cover ./...
.PHONY:mock-gen
mock-gen:
	mockgen -source=internal/adapters/app/service/person.go -destination=pkg/mocks/api/service/person_mock.go
	mockgen -source=internal/domain/ports/enricher/enricher.go -destination=pkg/mocks/api/enricher/enricher_mock.go
	mockgen -source=internal/domain/ports/repository/repository.go -destination=pkg/mocks/api/repository/repository_mock.go
.PHONY: migrate
migrate:
	goose -dir ./migrations postgres "postgres://admin:qwerty@localhost:5432/human?sslmode=disable" up
.PHONY: migrate-down
migrate-down:
	goose -dir ./migrations postgres "postgres://admin:qwerty@localhost:5432/human?sslmode=disable" down
# clean builds
.PHONY: clean
clean:
	go clean
	rm ${BINARY_NAME}-windows
	rm ${BINARY_NAME}-linux


