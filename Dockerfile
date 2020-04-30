FROM golang:latest

ENV PORT 8080
ENV DIMENSION 45
ENV ASSETS_PATH assets/default
ENV LIGHT_SQUARE_RGB 209,139,71
ENV DARK_SQUARE_RGB 255,206,158

WORKDIR /go/src/go-fen

ADD src src

RUN go build src/main.go

EXPOSE $PORT

CMD ["./main"]
