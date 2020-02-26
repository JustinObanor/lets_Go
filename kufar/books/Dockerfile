FROM golang:1.14.0-buster
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:3.11.3
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/main .
CMD ["./main"]