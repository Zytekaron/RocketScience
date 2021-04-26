package main

import (
	"RocketScience/src/events"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var token string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	token = "Bot " + os.Getenv("ROCKET_SCIENCE_TOKEN")
}

func main() {
	fmt.Println("Starting Rocket Science...")
	session, err := discordgo.New(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	session.AddHandler(events.MessageCreate)
	session.AddHandler(events.MessageEdit)

	err = session.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Connected to Discord!")

	fmt.Scanln()
}
