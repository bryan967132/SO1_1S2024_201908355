# detener contenedor por nombre
docker stop $(docker ps -aqf "name=grpc-server")
docker stop $(docker ps -aqf "name=grpc-client")
# eliminar contenedor por nombre
docker rm $(docker ps -aqf "name=grpc-server")
docker rm $(docker ps -aqf "name=grpc-client")
# eliminar imagen por nombre
docker rmi proyecto2_grpc_client
docker rmi proyecto2_grpc_server
# levantar docker-compose
docker-compose up -d