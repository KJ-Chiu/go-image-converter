FROM alpine:3.11.5 AS base
RUN apk --update add imagemagick

FROM golang:1.9.2 AS go-builder
WORKDIR /app
COPY main.go /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM base
COPY --from=go-builder /app/main /main
CMD ["/main"]
