FROM golang:latest AS builder
WORKDIR /backend
COPY . .
RUN go build -o Main .
EXPOSE 8000
CMD ["./Main"]