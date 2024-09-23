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

	gameState := gamelogic.NewGameState(name)

	err = pubsub.SubscribeJSON(conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+name,
		routing.PauseKey,
		pubsub.QueueTransient,
		handlerPause(gameState))
	if err != nil {
		log.Fatal(err)
	}

	armyMoveCh, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.ArmyMovesPrefix+".*",
		routing.ArmyMovesPrefix,
		pubsub.QueueTransient)
	if err != nil {
		log.Fatal(err)
	}

	err = pubsub.SubscribeJSON(conn,
		routing.ExchangePerilDirect,
		routing.ArmyMovesPrefix+"."+name,
		routing.ArmyMovesPrefix,
		pubsub.QueueTransient,
		handlerMove(gameState))
	if err != nil {
		log.Fatal(err)
	}

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
			// run spawn command
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// TODO publish
		} else if input[0] == "move" {
			// run move comand
			armyMove, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// publish move
			err = pubsub.PublishJSON(armyMoveCh, routing.ExchangePerilDirect, routing.ArmyMovesPrefix, armyMove)
			if err != nil {
				log.Fatalf("Failed to publish: %v \n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(p routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(p)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			return pubsub.Ack
		default:
			return pubsub.NackDiscard
		}
	}
}
