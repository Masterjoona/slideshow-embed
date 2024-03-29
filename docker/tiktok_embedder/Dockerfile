FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY *.go ./
ENV GOCACHE=/root/.cache/go-build
ARG build_tags=""
RUN --mount=type=cache,target="/root/.cache/go-build" go build -ldflags "-s -w" -tags=${build_tags}

FROM alpine:latest

RUN apk --no-cache add ffmpeg

WORKDIR /app

COPY --from=builder /app/meow /app/meow
COPY templates/ /app/templates

EXPOSE 4232
CMD ["./meow"]

