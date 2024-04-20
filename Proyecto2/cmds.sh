# go
go mod init P2
go mod tidy
go get -u google.golang.org/grpc
go get -u github.com/go-sql-driver/mysql

# protobuf
sudo apt update
sudo apt install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --go_out=. --go-grpc_out=. Client.proto
protoc --go_out=. --go-grpc_out=. Server.proto

# mysql
sudo docker-compose exec mysql mysql -uroot -p

#
go run Client/Client.go
go run Server/Server.go
locust -f Locust/Traffic.py

#
docker build -t dannyt8355/grpc_client .
docker run -p 3000:3000 dannyt8355/grpc_client
docker build -t dannyt8355/grpc_server .
docker run -p 3001:3001 dannyt8355/grpc_server