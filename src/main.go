package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
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
	case "why":
		session.ChannelMessageSend(msg.ChannelID, "Why not?")
	case "who":
		{
			var imageEmbed discordgo.MessageEmbedImage
			imageEmbed.URL = msg.Author.AvatarURL("128")

			var msgEmbed discordgo.MessageEmbedFooter
			msgEmbed.Text = msg.Member.JoinedAt.Local().Format(time.ANSIC)

			var embed discordgo.MessageEmbed
			embed.Title = msg.Author.Username
			embed.Image = &imageEmbed
			embed.Color = randColor()
			embed.Description = "Profile embed test"
			embed.Footer = &msgEmbed

			session.ChannelMessageSendEmbed(msg.ChannelID, &embed)
		}
	// change user profile picture with gelbooru/danbooru
	case "setpfp":
		{
			rand.Seed(time.Now().Unix())

			postNumber := rand.Intn(100) // hard cap for now
			url := fmt.Sprintf("https://testbooru.donmai.us/posts/%d", postNumber)
			fmt.Println(url)

			response, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
			}

			if response.StatusCode != 200 {
				fmt.Println("Bad request")
			}

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Failed to read body")
			}

			defer response.Body.Close()

			base := string(body)
			// find the image source (quick)
			// bad way to ignore mp4's, fix later
			if strings.Contains(base, ".mp4") {
				return
			}

			var lhs int
			var rhs int
			var bit int

			if strings.Contains(base, "sample") {
				bit = strings.Index(base, "sample")
			} else if strings.Contains(base, "original") {
				bit = strings.Index(base, "original")
			}

			lhs = strings.LastIndex(base[:bit], `"`)
			rhs = strings.Index(base[bit:], ">")
			// piece it together
			image := strings.Join(strings.Split(base[lhs:bit+rhs], `"`), "")

			fmt.Println(image)

			session.ChannelMessageSend(msg.ChannelID, string(image))
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

			session.ChannelMessageSend(msg.ChannelID, "cleared")
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
