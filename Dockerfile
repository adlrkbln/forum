FROM golang:1.22-alpine

WORKDIR /app


RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/web/main.go



RUN mkdir -p /data
VOLUME /data

EXPOSE 4000

CMD ["./main"]
