FROM golang:alpine AS builder
WORKDIR /tarea1Back
COPY . .
RUN go build -o Main .
EXPOSE 8000
CMD ["./Main"]