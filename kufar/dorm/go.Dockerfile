FROM golang:1.14.0-buster
WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y curl
RUN go build -o main .

FROM debian:10
WORKDIR /root/
COPY --from=0 /app/main .
EXPOSE 4200
CMD ["/root/main"]