docker start $(docker ps -aqf "name=grpc-server")
docker start $(docker ps -aqf "name=grpc-client")