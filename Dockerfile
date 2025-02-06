FROM golang:1.23-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o auth ./cmd/auth

FROM alpine:3.21.0

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build app/auth .
COPY --from=build app/.env .

CMD ["./auth"]
