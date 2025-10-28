FROM golang:1.25.1-alpine 

RUN apk add --no-cache git

WORKDIR /comments

COPY go.mod go.sum ./
COPY ./ ./

RUN go mod tidy

RUN go build -o comments ./cmd/main.go

CMD ["./comments"]