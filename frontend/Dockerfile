FROM golang:alpine

WORKDIR /usr/local/frontend

COPY . .

RUN go build -o frontend .

ENV HT_CONFIG_FILE /usr/local/frontend/config.yml

ENTRYPOINT ["./frontend"]