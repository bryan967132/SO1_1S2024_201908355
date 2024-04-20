go mod init T4
go mod tidy
go get -u google.golang.org/grpc
go get -u github.com/go-sql-driver/mysql

#protobuf
sudo apt update
sudo apt install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc -I . music.proto --go_out=. --go-grpc_out=.

#mysql
sudo docker-compose exec mysql mysql -uroot -p

#
go run Server/Server.go
locust -f Locust/Locust.py
go run Client/Client.go