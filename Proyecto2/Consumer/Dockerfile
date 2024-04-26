FROM golang:1.20-buster as build
WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v ./...
RUN go build -o /go/bin/app
FROM gcr.io/distroless/base-debian12
COPY --from=build /go/bin/app /
EXPOSE 3002
CMD ["/app"]