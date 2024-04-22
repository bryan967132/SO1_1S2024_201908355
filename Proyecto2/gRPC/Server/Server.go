package main

import (
	confproto "P2/Proto"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var ctx = context.Background()
var db *sql.DB

type Server struct {
	confproto.UnimplementedGetInfoServer
}

type Data struct {
	Name  string
	Album string
	Year  string
	Rank  string
}

func (d Data) toString() string {
	return fmt.Sprintf("Name: %v, Album: %v, Year: %v, Rank: %v", d.Name, d.Album, d.Year, d.Rank)
}

func kafka(data string) error {
	return nil
}

func (s *Server) ReturnInfo(ctx context.Context, in *confproto.RequestId) (*confproto.ReplyInfo, error) {
	fmt.Println("Recibí de cliente: ", in.GetName())
	data := Data{
		Name:  in.GetName(),
		Album: in.GetAlbum(),
		Year:  in.GetYear(),
		Rank:  in.GetRank(),
	}
	fmt.Println(data)
	err := kafka(data.toString())
	if err != nil {
		log.Printf("Failed to produce message to Kafka: %s", err)
	}
	return &confproto.ReplyInfo{Info: "Hola cliente, recibí el comentario"}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	confproto.RegisterGetInfoServer(s, &Server{})
	reflection.Register(s)

	if err := s.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}
