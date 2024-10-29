FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY pkg /app/pkg
COPY main.go /app/main.go
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -ldflags "-s -w"

FROM alpine:latest

RUN apk --no-cache add ffmpeg font-noto font-noto-cjk font-awesome font-noto-extra

WORKDIR /app

COPY --from=builder /app/meow /app/meow
COPY templates/ /app/templates

EXPOSE 4232
CMD ["./meow"]

