FROM golang:1.22-alpine

WORKDIR /app


COPY server.go .

RUN go build -o server server.go

ENV CHAT_PORT=1234

EXPOSE 1234

CMD ["./server"]
