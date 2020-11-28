FROM golang:alpine

WORKDIR /usr/local/backend

COPY . .

RUN go build -o backend .

ENTRYPOINT ["./backend"]