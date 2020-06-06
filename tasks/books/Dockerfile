FROM golang:1.14.0-buster
WORKDIR /app
COPY . .
RUN go build -o main .

FROM debian:10
WORKDIR /root/
COPY --from=0 /app/main .
ENV SESSION_KEY=default_session_key  
RUN apt-get update && apt-get install -y curl
CMD ["/root/main"]