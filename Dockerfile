FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go mod verify

COPY cmd/ ./cmd

COPY internal/ ./internal

RUN go build -o ./quotablegoofs ./cmd/web

EXPOSE 8000

ENTRYPOINT [ "./quotablegoofs" ]