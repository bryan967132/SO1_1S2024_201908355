FROM golang:1.22.2-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build Server.go
EXPOSE 3001
CMD ["./Server"]