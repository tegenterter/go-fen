FROM golang:latest AS builder

ADD src /app/
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

ARG DIMENSION
ARG ASSETS_PATH
ARG LIGHT_SQUARE_RGB
ARG DARK_SQUARE_RGB

ENV DIMENSION ${DIMENSION}
ENV ASSETS_PATH ${ASSETS_PATH}
ENV LIGHT_SQUARE_RGB ${LIGHT_SQUARE_RGB}
ENV DARK_SQUARE_RGB ${DARK_SQUARE_RGB}

WORKDIR /app/
COPY --from=builder /app .

EXPOSE 8080

CMD ["/app/main"]
