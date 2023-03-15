package clearcmd

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ClearCommand struct{}

func New() *ClearCommand {
	return &ClearCommand{}
}

func (*ClearCommand) Name() string {
	return "clear"
}

func (*ClearCommand) Description() string {
	return "Purges all messages in a channel"
}

func (*ClearCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	var msgIDs []string
	var oldMsgIDs []*discordgo.Message

	msgs, _ := s.ChannelMessages(ic.ChannelID, 100, "", "", "")
	for _, msg := range msgs {
		if time.Since(msg.Timestamp).Hours() > 336 {
			oldMsgIDs = append(oldMsgIDs, msg)
		} else {
			msgIDs = append(msgIDs, msg.ID)
		}
	}

	err := s.ChannelMessagesBulkDelete(ic.ChannelID, msgIDs)
	if err != nil {
		log.Println("Failed to delete all messages in this channel:", err)
	}

	for _, msg := range oldMsgIDs {
		err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if err != nil {
			log.Println("Failed to delete message:", err)
			break
		}
	}

	s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Cleared!",
		},
	})
}
