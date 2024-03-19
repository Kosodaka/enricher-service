FROM golang:1.21 as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/enricher-service ./cmd/app/main.go

FROM alpine:3.18.6
COPY --from=builder /build/bin/enricher-service /
COPY migrations /
COPY .env /
ENTRYPOINT ["/enricher-service"]
