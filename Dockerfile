FROM golang:1.22-alpine

WORKDIR /app

RUN apk update && apk add --no-cache make git bash && apk add curl

RUN chmod 755 /app

COPY go.mod go.sum ./
RUN go mod download

CMD ["go", "run", "main.go"]
