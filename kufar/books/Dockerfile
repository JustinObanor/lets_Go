FROM golang:1.14.0-buster
WORKDIR /app
COPY . .
RUN go build -o main .

FROM debian:10
WORKDIR /root/
COPY --from=0 /app/main .
# ENV SESSION_KEY $SESSION_KEY 
# RUN echo $SESSION_KEY
CMD ["/root/main"]