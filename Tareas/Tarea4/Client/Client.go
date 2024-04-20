package main

import (
	"context"
	"log"
	"os"

	music "T4/Protobuf"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 5 {
		log.Fatalf("Usage: %s <name> <album> <year> <rank>", os.Args[0])
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := music.NewMusicServiceClient(conn)

	name := os.Args[0]
	album := os.Args[1]
	year := os.Args[2]
	rank := os.Args[3]

	r, err := c.SendMusicInfo(context.Background(), &music.MusicRequest{Name: name, Album: album, Year: year, Rank: rank})
	if err != nil {
		log.Fatalf("could not add band: %v", err)
	}
	log.Printf("Response: %s", r.Message)
}
