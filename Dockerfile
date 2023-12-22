FROM golang:1.17-alpine as build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o services_cli ./cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=build /app/services_cli .
ENTRYPOINT ["./services_cli", "-s"]
CMD ["-s", ".:/app"]
