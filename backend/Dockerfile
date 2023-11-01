FROM golang:1.21-alpine AS backend
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download && go mod verify
COPY . .
RUN go build -o /app/bar /app/cmd/bar/main.go

FROM alpine
COPY --from=backend /app/bar .
USER 1000
ENTRYPOINT ["/bar"]
