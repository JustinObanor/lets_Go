FROM golang:1.14.0-buster

WORKDIR /opt/multiplexer

COPY . .

RUN go build 

FROM debian:10

RUN apt-get update && \
    apt-get install -y ca-certificates

WORKDIR /opt/multiplexer

COPY --from=0 /opt/multiplexer/multiplexer1 .

EXPOSE 8080

CMD [ "./multiplexer1" ]