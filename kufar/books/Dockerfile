FROM golang:1.14.0-buster
WORKDIR /app
COPY . .
RUN go build -o main .

FROM debian:10
WORKDIR /root/
COPY --from=0 /app/main .
ENV SESSION_AUTH_KEY=default_session_key  
ENV SESSION_ENCRYPT_KEY=default_session_key
CMD ["/root/main"]