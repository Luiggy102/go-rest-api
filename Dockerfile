FROM golang:1.22.7-alpine AS constructor

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o rest-api

FROM alpine:latest

WORKDIR /app
COPY . .
COPY --from=constructor /app/rest-api .
CMD ["./rest-api"]
