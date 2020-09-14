FROM golang:latest
ADD . /server
WORKDIR /server
RUN go mod download
RUN go build server.go
EXPOSE 5001
CMD ["./server"]