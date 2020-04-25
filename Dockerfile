FROM golang:1.13

RUN mkdir /go-fen
ADD . /go-fen/
WORKDIR /go-fen

RUN go build -o main .

CMD ["/go-fen/main"]
