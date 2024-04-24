FROM rust:latest as builder
WORKDIR /usr/src/app
COPY . .
RUN cargo build --release
FROM ubuntu:latest
RUN apt-get update && apt-get install -y libssl-dev
COPY --from=builder /usr/src/app/target/release/Client /usr/local/bin/Client
EXPOSE 4000
CMD ["Client"]