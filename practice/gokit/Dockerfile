FROM golang:1.14.0-buster
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

FROM debian:10
WORKDIR /root/
COPY --from=0 /app/main .
EXPOSE 8080
CMD ["/root/main"]