FROM golang:1.14.6-alpine3.12
RUN mkdir /dione
COPY . /dione
WORKDIR /dione

RUN apk add git
RUN apk add --update make
RUN go mod download

RUN make build

CMD ["./dione"]