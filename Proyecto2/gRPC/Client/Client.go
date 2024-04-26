package main

import (
	confproto "P2/Proto"
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ctx = context.Background()

type Data struct {
	Name  string
	Album string
	Year  string
	Rank  string
}

func insertData(c *fiber.Ctx) error {
	var data map[string]string
	e := c.BodyParser(&data)
	if e != nil {
		return e
	}

	voto := Data{
		Name:  data["name"],
		Album: data["album"],
		Year:  data["year"],
		Rank:  data["rank"],
	}
	fmt.Println(voto)
	go sendServer(voto)
	return nil
}

func sendServer(voto Data) {
	conn, err := grpc.Dial("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}

	cl := confproto.NewGetInfoClient(conn)
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(conn)

	ret, err := cl.ReturnInfo(ctx, &confproto.RequestId{
		Name:  voto.Name,
		Album: voto.Album,
		Year:  voto.Year,
		Rank:  voto.Rank,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ret.GetInfo())
}

func main() {
	app := fiber.New()

	app.Get("/grpc", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "Client is running!!!",
		})
	})
	app.Post("/grpc/insert", insertData)

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
