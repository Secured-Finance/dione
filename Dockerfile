FROM golang:1.14.6-alpine3.12
RUN mkdir /p2p-oracle-node
COPY . /p2p-oracle-node
WORKDIR /p2p-oracle-node

RUN apk add git
RUN apk add --update make
RUN go mod download

RUN make build

CMD ["./main"]