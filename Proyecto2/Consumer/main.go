package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	Name  string `json:"Name"`
	Album string `json:"Album"`
	Year  string `json:"Year"`
	Rank  string `json:"Rank"`
}

func getRedisClient() (context.Context, *redis.Client) {
	var ctx context.Context = context.Background()
	var client *redis.Client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "YOUR_PASSWORD",
		DB:       0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error al hacer ping a Redis:", err)
		return nil, nil
	}
	fmt.Println("Respuesta de Redis:", pong)
	return ctx, client
}

func getMongoConnection() (context.Context, *mongo.Collection) {
	var ctx context.Context = context.Background()
	var client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %s", err)
		return nil, nil
	}
	defer client.Disconnect(ctx)

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Error al hacer ping a MongoDB:", err)
		return nil, nil
	}
	fmt.Println("Conexión exitosa a MongoDB. Ping OK.")

	var logCollection *mongo.Collection = client.Database("bdso1p2").Collection("logs")
	return ctx, logCollection
}

func convertData(data string) *Data {
	var newData Data
	err := json.Unmarshal([]byte(data), &newData)
	if err != nil {
		fmt.Println("Error al convertir la cadena a JSON:", err)
		return nil
	}
	return &newData
}

func insertRedis(ctx context.Context, rdb *redis.Client, data *Data) error {
	hashKey := fmt.Sprintf("%s:%s:%s", data.Name, data.Album, data.Year)
	fmt.Printf("hashkey new vote = %v\n", hashKey)

	err := rdb.HIncrBy(ctx, "albums", hashKey, 1).Err()
	if err != nil {
		return err
	}

	fmt.Printf("Redis -> Se agregó album: %v, hashkey: %v\n", data.Album, hashKey)

	_, err = rdb.Incr(ctx, "total").Result()
	if err != nil {
		return err
	}

	fmt.Println("Redis -> Se incrementó el total de votos")

	return nil
}

func insertMongo(ctx context.Context, collection *mongo.Collection, data *Data) error {
	log := bson.M{
		"Name": data.Name, "ALbum": data.Album, "Year": data.Year, "Rank": data.Rank,
	}

	_, err := collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}

	fmt.Printf("Mongo -> Log insertado en MongoDB: %v\n", data)
	return nil
}

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-bootstrap:9092",
		"group.id":          "mygroupid",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create Kafka consumer: %s", err)
		os.Exit(1)
	}

	var redisCtx, redisClient = getRedisClient()

	var mongoCtx, mongoCollection = getMongoConnection()

	topic := "topic-so1p2"
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s", err)
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

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

				fmt.Printf("Consumed event from topic %s: value = %s\n", *ev.TopicPartition.Topic, string(ev.Value))

				err = insertRedis(redisCtx, redisClient, convertData(string(ev.Value)))
				if err != nil {
					fmt.Printf("Failed to process and update Redis: %s\n", err)
				}

				err = insertMongo(mongoCtx, mongoCollection, convertData(string(ev.Value)))
				if err != nil {
					fmt.Printf("Failed to insert log into MongoDB: %s\n", err)
				}

			}
		}
	}()

	<-sigchan
}
