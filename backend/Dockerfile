FROM golang:1.21-alpine AS backend
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download && go mod verify
COPY . .
RUN go build -o /app/bar /app/cmd/bar/main.go

FROM golang:1.21-alpine AS nfc
RUN apk add gcc libc-dev pcsc-lite-dev
WORKDIR /app
ENV CGO_ENABLED=1
COPY nfc/go.mod nfc/go.sum ./
COPY nfc .
RUN go build -o /app/nfc /app/*.go

FROM alpine
COPY --from=backend /app/bar .
COPY --from=nfc /app/nfc .
USER 1000
ENTRYPOINT ["/bar"]
