FROM golang:1.14-alpine as build

WORKDIR /usr/src/app

RUN mkdir bin
COPY . .
RUN go mod download
RUN go build -o bin/app main.go


FROM alpine:latest
RUN apk add ca-certificates

COPY cert.pem /usr/local/bin/certificates/cert.pem
COPY private.key /usr/local/bin/certificates/private.key

COPY --from=build /usr/src/app/bin/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]