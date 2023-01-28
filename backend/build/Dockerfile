FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o gofile-upload.bin cmd/main.go

EXPOSE 8080
CMD ["./gofile-upload.bin"]