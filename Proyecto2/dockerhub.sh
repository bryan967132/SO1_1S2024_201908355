# docker rmi dannyt8355/grpc-client:p1
# docker rmi dannyt8355/grpc-server:p1
# cd gRPC/Client
# docker build --no-cache -t dannyt8355/grpc-client:p1 .
# cd ../Server
# docker build --no-cache -t dannyt8355/grpc-server:p1 .
# cd ../../
# docker push dannyt8355/grpc-client:p1
# docker push dannyt8355/grpc-server:p1
# =========================================================
# docker rmi dannyt8355/rust-client:p1
# docker rmi dannyt8355/rust-server:p1
# cd Rust/Client
# docker build --no-cache -t dannyt8355/rust-client:p1 .
# cd ../Server
# docker build --no-cache -t dannyt8355/rust-server:p1 .
# cd ../../
# docker push dannyt8355/rust-client:p1
# docker push dannyt8355/rust-server:p1
# =========================================================
# docker rmi dannyt8355/consumer:p1
cd Consumer
docker build --no-cache -t dannyt8355/consumer:latest .
cd ../
docker push dannyt8355/consumer:latest