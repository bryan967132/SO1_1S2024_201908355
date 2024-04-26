package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	Album  string `json:"Album"`
	Year   string `json:"Year"`
	Artist string `json:"Artist"`
	Ranked string `json:"Ranked"`
}

type Log struct {
	Data string `bson:"data"`
}

func main() {
	// Configuración de Kafka
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-bootstrap:9092",
		"group.id":          "mygroupid",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create Kafka consumer: %s", err)
		os.Exit(1)
	}

	// Configuración de Redis
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "YOUR_PASSWORD",
		DB:       0,
	})

	// Realiza un ping a Redis
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error al hacer ping a Redis:", err)
		return
	}
	fmt.Println("Respuesta de Redis:", pong)

	// Configuración de MongoDB
	mongoURI := "mongodb://mongodb:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Failed to create MongoDB client: %s", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %s", err)
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Error al hacer ping a MongoDB:", err)
		return
	}
	fmt.Println("Conexión exitosa a MongoDB. Ping OK.")

	database := client.Database("bdso1p2")
	logCollection := database.Collection("logs")

	topic := "topic-so1p2"
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s", err)
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Función para consumir mensajes de Kafka,  actualizar Redis y insertar mongo
	go func() {
		for {
			select {
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				c.Close()
				os.Exit(0)
			default:

				ev, err := c.ReadMessage(100 * time.Millisecond)
				if err != nil {
					continue
				}
				fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
					*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))

				// Actualizar datos en Redis
				err = processAndUpdateRedis(ctx, rdb, string(ev.Value))
				if err != nil {
					fmt.Printf("Failed to process and update Redis: %s\n", err)
				}

				// Insertar log en MongoDB
				err = insertLog(ctx, logCollection, string(ev.Value))
				if err != nil {
					fmt.Printf("Failed to insert log into MongoDB: %s\n", err)
				}

			}
		}
	}()

	// Mantener el programa en ejecución
	<-sigchan
}

func processAndUpdateRedis(ctx context.Context, rdb *redis.Client, data string) error {
	// Procesar la cadena de datos para extraer los valores
	values := strings.Split(data, ", ")
	name := strings.Split(values[0], ": ")[1]
	album := strings.Split(values[1], ": ")[1]
	year := strings.Split(values[2], ": ")[1]

	// Clave del hash en Redis
	hashKey := fmt.Sprintf("albums:%s:%s:%s", name, album, year)
	fmt.Printf("hashkey new vote = %v\n", hashKey)

	// Actualizar el contador de votos en Redis
	err := rdb.HIncrBy(ctx, "albums", hashKey, 1).Err()
	if err != nil {
		return err
	}

	fmt.Printf("Se agregó albums:" + album + " hashkey: " + hashKey)
	return nil
}

func insertLog(ctx context.Context, collection *mongo.Collection, data string) error {
	// Procesar la cadena de datos para extraer los valores
	values := strings.Split(data, ", ")
	name := strings.Split(values[0], ": ")[1]
	album := strings.Split(values[1], ": ")[1]
	year := strings.Split(values[2], ": ")[1]
	rank := strings.Split(values[3], ": ")[1]
	data = "ALbum: " + album + " Name: " + name + " year: " + year + " rank: " + rank
	log := Log{
		Data: data,
	}

	_, err := collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}

	fmt.Println("Log insertado en MongoDB: " + data)
	return nil
}
