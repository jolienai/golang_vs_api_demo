FROM golang:1.22-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app

COPY --from=build /out/api /api
COPY migrations /migrations

USER app
EXPOSE 8080

ENTRYPOINT ["/api"]
