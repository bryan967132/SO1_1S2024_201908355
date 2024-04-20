docker stop $(docker ps -aqf "name=grpc-server")
docker stop $(docker ps -aqf "name=grpc-client")