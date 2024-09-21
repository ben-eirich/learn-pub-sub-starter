package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	con, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ. Is it running?\n %v\n", err)
	}
	defer con.Close()
	fmt.Println("RabbitMQ connection successful!")

	mqch, err := con.Channel()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		} else if input[0] == "quit" {
			fmt.Println("Exiting")
			break
		} else if input[0] == "help" {
			gamelogic.PrintServerHelp()
		} else if input[0] == "pause" {
			err = pubsub.PublishJSON(mqch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Fatalf("Failed to publish: %v \n", err)
			}
		} else if input[0] == "resume" {
			err = pubsub.PublishJSON(mqch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Fatalf("Failed to publish: %v \n", err)
			}
		} else {
			fmt.Printf("I didn't understand '%s' command.\n", input[0])
		}
	}

	fmt.Println("Shutting down.")
}
