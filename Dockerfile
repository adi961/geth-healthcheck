FROM golang:1.16.3 as build
WORKDIR /usr/src/app

COPY . .
RUN go mod download -x
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v ./cmd/healthcheck/

FROM alpine:latest

COPY --from=build /usr/src/app/healthcheck /

CMD ["/healthcheck"]
