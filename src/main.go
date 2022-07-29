package main

import (
	//"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	//"io/ioutil"
	//"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	//args := os.Args[1]
	//token := args
	token := os.Getenv("DAUTH_TOKEN")

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

func randColor() int {
	rand.Seed(time.Now().Unix())

	colors := []int{
		0xFF1567,
		0x9B74FF,
		0xFFDDE4,
	}

	return colors[rand.Intn(len(colors))]
}

func createMsg(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if msg.Author.ID == session.State.User.ID {
		return
	}

	switch msg.Content {
	case "ping":
		session.ChannelMessageSend(msg.ChannelID, "Pong!")
	case "pong":
		session.ChannelMessageSend(msg.ChannelID, "Ping!")
	case "hello":
		session.ChannelMessageSend(msg.ChannelID, "I'm living in your walls")
	case "why?":
		session.ChannelMessageSend(msg.ChannelID, "Why not?")
	case "who am I":
		{
			var imageEmbed discordgo.MessageEmbedImage
			imageEmbed.URL = msg.Author.AvatarURL("128")

			var embed discordgo.MessageEmbed
			embed.Title = msg.Author.Username
			embed.Image = &imageEmbed
			embed.Color = randColor()
			embed.Description = "Profile embed test"

			session.ChannelMessageSendEmbed(msg.ChannelID, &embed)
		}
		// purge messages in current channel
	case "cls":
		{
			var msgArr []string
			msgs, _ := session.ChannelMessages(msg.ChannelID, 100, "", "", "")
			for _, m := range msgs {
				msgArr = append(msgArr, m.ID)
			}

			err := session.ChannelMessagesBulkDelete(msg.ChannelID, msgArr)
			if err != nil {
				fmt.Println("Mini failed to delete all messages in all channels")
			}
		}
	// purge messages in all channels
	case "cls-all":
		{
			var msgArr []string
			// loop through all channels
			channels, _ := session.GuildChannels(msg.GuildID)
			for _, curChannel := range channels {
				if curChannel.Type != discordgo.ChannelTypeGuildText {
					continue
				}

				msgs, _ := session.ChannelMessages(curChannel.ID, 100, "", "", "")
				for _, m := range msgs {
					msgArr = append(msgArr, m.ID)
				}

				err := session.ChannelMessagesBulkDelete(curChannel.ID, msgArr)
				if err != nil {
					fmt.Println("Mini failed to delete all messages in all channels")
				}

				session.ChannelMessageSend(curChannel.ID, "cleared")
			}
		}
	}
}
