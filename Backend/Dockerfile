FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o backend .

RUN apk --no-cache add ca-certificates

CMD ["./backend"]