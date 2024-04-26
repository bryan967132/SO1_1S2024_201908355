FROM golang:latest AS builder
WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v ./...
RUN go build -o /go/bin/app
FROM gcr.io/distroless/base-debian12
COPY --from=builder /go/bin/app /
CMD ["/app"]