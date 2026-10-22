FROM golang:1.27-alpine AS build
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download && go mod verify
COPY . .
RUN go build -o /app/bar /app/cmd/bar/main.go

FROM alpine:3.24 AS runtime
USER 1000
ENTRYPOINT ["/bar"]

# Release image: GoReleaser builds the binary and puts it in the context
FROM runtime AS release
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/bar /bar

# Default target, built from source for docker compose
FROM runtime
COPY --from=build /app/bar /bar
