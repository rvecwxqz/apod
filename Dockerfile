FROM golang:1.21.5

WORKDIR /

COPY . .

RUN go build -o ./app ./cmd/app/

CMD ["./app"]