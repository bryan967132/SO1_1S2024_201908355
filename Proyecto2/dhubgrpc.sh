docker rmi dannyt8355/grpc-client:latest
docker rmi dannyt8355/grpc-server:latest
cd gRPC/Client
docker build -t dannyt8355/grpc-client:latest .
cd ../Server
docker build -t dannyt8355/grpc-server:latest .
cd ../
docker push dannyt8355/grpc-client:latest
docker push dannyt8355/grpc-server:latest