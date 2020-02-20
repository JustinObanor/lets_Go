FROM golang:latest
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq
RUN go build -o main .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
CMD ["/app/main"]