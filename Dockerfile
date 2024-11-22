FROM golang:1.22-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum main.go ./
RUN go mod download
RUN go mod tidy

COPY api ./api
COPY common ./common
COPY config ./config
COPY consumer ./consumer
COPY db ./db
COPY models ./models

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

FROM scratch

COPY --from=builder /dist/main .

COPY --from=builder /dist/*.properties .

ENV PROFILE=prod \
    DATABASE_HOST=${DATABASE_HOST} \
    DATABASE_NAME=${DATABASE_NAME} \
    DATABASE_USER=${DATABASE_USER} \
    DATABASE_PASSWORD=${DATABASE_PASSWORD} \
    DATABASE_PORT=${DATABASE_PORT}

ENTRYPOINT ["/main"]
