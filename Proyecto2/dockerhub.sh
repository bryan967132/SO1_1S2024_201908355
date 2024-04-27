docker rmi dannyt8355/grpc-client:grpc
docker rmi dannyt8355/grpc-server:grpc
cd gRPC/Client
docker build --no-cache -t dannyt8355/grpc-client:grpc .
cd ../Server
docker build --no-cache -t dannyt8355/grpc-server:grpc .
cd ../../
docker push dannyt8355/grpc-client:grpc
docker push dannyt8355/grpc-server:grpc
# =========================================================
docker rmi dannyt8355/rust-client:rust
docker rmi dannyt8355/rust-server:rust
cd Rust/Client
docker build --no-cache -t dannyt8355/rust-client:rust .
cd ../Server
docker build --no-cache -t dannyt8355/rust-server:rust .
cd ../../
docker push dannyt8355/rust-client:rust
docker push dannyt8355/rust-server:rust
# # =========================================================
docker rmi dannyt8355/kafka-consumer:latest
cd Consumer
docker build --no-cache -t dannyt8355/kafka-consumer:latest .
cd ../
docker push dannyt8355/kafka-consumer:latest