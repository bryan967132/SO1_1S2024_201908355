FROM rust:latest as builder
WORKDIR /usr/src/app
COPY . .
RUN cargo build --release
FROM ubuntu:latest
RUN apt-get update && apt-get install -y libssl-dev
COPY --from=builder /usr/src/app/target/release/Server /usr/local/bin/Server
EXPOSE 4001
CMD ["Server"]