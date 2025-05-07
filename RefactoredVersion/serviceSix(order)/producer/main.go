package main

import (
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	userId := "124"
	message := "You have a new order!"
	err = nc.Publish("order.notification."+userId, []byte(message))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Message sent to", userId)
}
