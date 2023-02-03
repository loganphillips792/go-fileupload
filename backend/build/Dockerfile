FROM golang:latest

WORKDIR /app
RUN ls
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o gofile-upload.bin cmd/main.go

EXPOSE 8000
CMD ["./gofile-upload.bin"]