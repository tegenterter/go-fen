FROM golang:latest

RUN mkdir /go-fen
ADD . /go-fen/
WORKDIR /go-fen

RUN go build -o main .

CMD ["/go-fen/main"]
