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
	fmt.Println("Starting Peril client...")

	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ. Is it running?\n %v\n", err)
	}
	defer conn.Close()
	fmt.Println("RabbitMQ connection successful!")

	name, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", name),
		routing.PauseKey,
		pubsub.QueueTransient)
	if err != nil {
		log.Fatal(err)
	}

	gameState := gamelogic.NewGameState(name)
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		} else if input[0] == "quit" {
			gamelogic.PrintQuit()
			return
		} else if input[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if input[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		} else if input[0] == "status" {
			gameState.CommandStatus()
		} else if input[0] == "spawn" {
			gameState.CommandSpawn(input)
			// TODO check error, and publish
		} else if input[0] == "move" {
			gameState.CommandMove(input)
			// TODO check error, and publish
		} else {
			fmt.Println("Unknown command")
		}
	}
}
