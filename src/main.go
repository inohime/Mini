package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// for now, just make all commands in one file
var (
	errSymbol = "[X]"

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "carrot",
			Description: "Says carrot in the current channel!",
		},
		{
			Name:        "cls",
			Description: "Purges all messages in a channel",
		},
		{
			Name:        "cls-all",
			Description: "Purges all messages in every channels",
		},
		{
			Name:        "generate",
			Description: "Generates a new profile picture",
		},
	}

	commandList = map[string]func(session *discordgo.Session, ic *discordgo.InteractionCreate){
		"carrot": func(session *discordgo.Session, ic *discordgo.InteractionCreate) {
			session.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Carrot!",
				},
			})
		},
		"cls": func(session *discordgo.Session, ic *discordgo.InteractionCreate) {
			var msgArr []string
			today := time.Now()
			msgs, _ := session.ChannelMessages(ic.ChannelID, 100, "", "", "")
			for _, m := range msgs {
				// slow, fix later
				if !inTimeSpan(m.Timestamp, m.Timestamp.AddDate(0, 0, 14), today) {
					err := session.ChannelMessageDelete(m.ChannelID, m.ID)
					if err != nil {
						log.Println(errSymbol+" Failed to delete message:", err)
					}
				} else {
					msgArr = append(msgArr, m.ID)
				}
			}

			err := session.ChannelMessagesBulkDelete(ic.ChannelID, msgArr)
			if err != nil {
				log.Println(errSymbol+" Mini failed to delete all messages in this channel:", err)
			}

			session.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "cleared",
				},
			})

			defer time.AfterFunc(time.Second*10, func() {
				session.InteractionResponseDelete(ic.Interaction)
			})
		},
		"cls-all": func(session *discordgo.Session, ic *discordgo.InteractionCreate) {
			var msgArr []string
			today := time.Now()
			// loop through all channels
			channels, _ := session.GuildChannels(ic.GuildID)
			for _, curChannel := range channels {
				if curChannel.Type != discordgo.ChannelTypeGuildText {
					continue
				}

				// check against array and not this..
				if curChannel.ID == "1040744416001937531" || curChannel.ID == "1002590272531726446" || curChannel.ID == "1040744416001937530" {
					// skip mod, cheese, rules
					continue
				}

				log.Println("Clearing Channel:", curChannel.Name)

				msgs, _ := session.ChannelMessages(curChannel.ID, 100, "", "", "")
				for _, m := range msgs {
					// slow, fix later
					if !inTimeSpan(m.Timestamp, m.Timestamp.AddDate(0, 0, 14), today) {
						err := session.ChannelMessageDelete(m.ChannelID, m.ID)
						if err != nil {
							log.Println(errSymbol+" Failed to delete message:", err)
						}
					} else {
						msgArr = append(msgArr, m.ID)
					}
				}

				err := session.ChannelMessagesBulkDelete(curChannel.ID, msgArr)
				if err != nil {
					log.Println(errSymbol+" Mini failed to delete all messages in this channel:", err)
				}

				session.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "cleared",
					},
				})
			}
		},
		"generate": func(session *discordgo.Session, ic *discordgo.InteractionCreate) {
			rand.Seed(time.Now().Unix())

			postNumber := rand.Intn(100) // hard cap for now
			url := fmt.Sprintf("https://testbooru.donmai.us/posts/%d", postNumber)
			log.Println("Image Requested:", url)

			response, err := http.Get(url)
			if err != nil {
				log.Println(errSymbol, err)
				return
			}

			if response.StatusCode != 200 {
				log.Println(errSymbol+" Bad request:", err)
				return
			}

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Println(errSymbol+" Failed to read body:", err)
				return
			}

			defer response.Body.Close()

			base := string(body)

			log.Println("Body Requested:", base)
			// find the image source (quick)
			// bad way to ignore mp4's, fix later
			if strings.Contains(base, ".mp4") {
				return
			}

			var bit int
			if strings.Contains(base, "sample") {
				bit = strings.Index(base, "sample")
			} else if strings.Contains(base, "original") {
				bit = strings.Index(base, "original")
			}

			lhs := strings.LastIndex(base[:bit], `"`)
			rhs := strings.Index(base[bit:], ">")
			// piece it together
			image := strings.Join(strings.Split(base[lhs:bit+rhs], `"`), "")

			// find the artist

			// find the image tags (general, copyrights, artist, characters)

			var footerEmbed discordgo.MessageEmbedFooter
			footerEmbed.Text = "Requested at: " + time.Now().Local().Format(time.ANSIC)

			// create an embed with the data
			var imageEmbed discordgo.MessageEmbedImage
			// modify image size for embeds
			imageEmbed.URL = image

			session.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Recommendation for " + ic.Member.User.Username,
							Image:       &imageEmbed,
							Color:       randColor(),
							Description: "Profile embed test",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Artist",
									Value: "Artist Name here",
								},
								{
									Name:  "Tags",
									Value: "Tags here",
								},
								{
									Name:  "Link",
									Value: "Image link here",
								},
							},
							Footer: &footerEmbed,
						},
					},
					AllowedMentions: &discordgo.MessageAllowedMentions{},
				},
			})
		},
	}
)

// small helper functions
func randColor() int {
	rand.Seed(time.Now().Unix())

	colors := []int{
		0xFF1567,
		0x9B74FF,
		0xFFDDE4,
	}

	return colors[rand.Intn(len(colors))]
}

func inTimeSpan(startDate, endDate, valid time.Time) bool {
	return valid.After(startDate) && valid.Before(endDate)
}

func main() {
	token := os.Getenv("DAUTH_TOKEN")

	newSession, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println(errSymbol+" Error waking up Mini:", err)
		return
	}
	// verify the correct token was set successfully
	//log.Println(newSession.Identify.Token)

	newSession.AddHandler(func(session *discordgo.Session, ic *discordgo.InteractionCreate) {
		if cmd, ok := commandList[ic.ApplicationCommandData().Name]; ok {
			cmd(session, ic)
		}
	})

	newSession.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		log.Printf("Logged in as %v#%v", session.State.User.Username, session.State.User.Discriminator)
	})

	err = newSession.Open()
	if err != nil {
		log.Println(errSymbol+" Error opening socket:", err)
		return
	}

	// add our new commands
	log.Println("Adding commands...")
	addedCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for ic, cmd := range commands {
		makeCommand, err := newSession.ApplicationCommandCreate(newSession.State.User.ID, "", cmd)
		if err != nil {
			log.Panicf(errSymbol+" Failed to create [%v] command: %v", cmd.Name, err)
		}
		addedCommands[ic] = makeCommand
	}

	log.Println("Mini is awake, press CTRL+C to sleep")

	sigInput := make(chan os.Signal, 1)
	signal.Notify(sigInput, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sigInput

	// pop commands
	log.Println("Removing commands...")
	for _, cmd := range addedCommands {
		err := newSession.ApplicationCommandDelete(newSession.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Panicf(errSymbol+" Failed to delete [%v] command: %v", cmd.Name, err)
		}
	}

	// close the connection
	newSession.Close()
}
