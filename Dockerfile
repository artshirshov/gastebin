FROM golang:alpine

WORKDIR /gastebin
COPY . .

RUN go build -o ./bin/pastebin ./cmd/pastebin

CMD ["/gastebin/bin/pastebin"]
EXPOSE 8080