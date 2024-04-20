package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	music "T4/Protobuf"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	music.UnimplementedMusicServiceServer
}

func (s *server) SendMusicInfo(ctx context.Context, req *music.MusicRequest) (*music.MusicResponse, error) {
	db, err := sql.Open("mysql", "root:mysqlpass@tcp(localhost:3306)/music_data")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO music (name, album, year, ranking) VALUES (?, ?, ?, ?)", req.Name, req.Album, req.Year, req.Rank)
	if err != nil {
		return nil, err
	}

	return &music.MusicResponse{Message: "Data received and stored"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	music.RegisterMusicServiceServer(s, &server{})
	reflection.Register(s)
	fmt.Println("Server started at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
