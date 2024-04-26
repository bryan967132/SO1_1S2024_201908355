package main

import (
	"context"
	"fmt"
	"log"
	"net"

	confproto "P2/Proto"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/grpc"
)

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

func producerkafka(data string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-bootstrap:9092",
	})
	if err != nil {
		return err
	}
	defer p.Close()

	fmt.Println("PRODUCER")

	deliveryChan := make(chan kafka.Event)

	var kafkaTopic = "topic-so1p2"
	_ = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          []byte(data),
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

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
	err := producerkafka(data.toString())
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

	if err := s.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}
