package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	args := os.Args[1]
	token := args

	newSession, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error waking up Mini", err)
		return
	}
	// verify the correct token was set successfully
	fmt.Println(newSession.Identify.Token)

	newSession.AddHandler(createMsg)

	newSession.Identify.Intents = discordgo.IntentsGuildMessages

	err = newSession.Open()
	if err != nil {
		fmt.Println("Error opening socket", err)
		return
	}

	fmt.Println("Mini is awake, press CTRL-C to sleep")
	signalInput := make(chan os.Signal, 1)

	signal.Notify(signalInput, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-signalInput

	// close the connection
	newSession.Close()
}

func createMsg(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if msg.Author.ID == session.State.User.ID {
		return
	}

	if msg.Content == "ping" {
		session.ChannelMessageSend(msg.ChannelID, "Pong!")
	}

	if msg.Content == "pong" {
		session.ChannelMessageSend(msg.ChannelID, "Ping!")
	}

	if msg.Content == "hello" {
		session.ChannelMessageSend(msg.ChannelID, "I'm living in your walls")
	}
}
