FROM golang:latest as build

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY pkg pkg
COPY internal internal
COPY cmd/api/main.go cmd/api/main.go
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /api indexstorm/go-api-boilerplate/cmd/api

FROM gcr.io/distroless/static

WORKDIR /

COPY --from=build /api /api
COPY --chown=nonroot:nonroot internal/jwt /etc/jwt
COPY --chown=nonroot:nonroot internal/pg.pem /etc/pg.pem

USER nonroot:nonroot
ENTRYPOINT ["/api"]
