FROM golang:1.14.0-buster
WORKDIR /app
COPY . .
RUN go build ./cmd/student/main.go

FROM debian:10
WORKDIR /app
COPY --from=0 /app/main .
EXPOSE 8080
CMD ["/app/main"]