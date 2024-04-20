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

func mysqlConnect() {
	var err error
	db, err = sql.Open("mysql", "root:mysqlpass@tcp(mysql:3306)/P2SO1")
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Conexión a MySQL exitosa")
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
	query := "INSERT INTO votos (name, album, year, ranking) VALUES (?, ?, ?, ?)"
	_, err := db.ExecContext(ctx, query, data.Name, data.Album, data.Year, data.Rank)
	if err != nil {
		log.Println("Error al insertar en MySQL:", err)
	}
	return &confproto.ReplyInfo{Info: "Hola cliente, recibí el comentario"}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	mysqlConnect()
	confproto.RegisterGetInfoServer(s, &Server{})
	reflection.Register(s)

	if err := s.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}
