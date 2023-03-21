package helpcmd

import (
	"fmt"
	base "main/src/ops"
	"main/src/utils"
	"time"

	"github.com/bwmarrin/discordgo"
)

type HelpCommand struct{}

func New() *HelpCommand {
	return &HelpCommand{}
}

func (*HelpCommand) Name() string {
	return "help"
}

func (*HelpCommand) Description() string {
	return "Purges all messages in every channel"
}

func (c *HelpCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Table of Commands",
					Color: utils.RandomColor(),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name: "carrot",
							Value: fmt.Sprintf(
								"Creates a blob with the word 'carrot' in it.\n Usage: %s",
								"/carrot",
							),
						},
						{
							Name: "clear",
							Value: fmt.Sprintf(
								"Give it a channel and it'll purge all messages.\n Usage: %s",
								"/clear <channel>",
							),
						},
						{
							Name: "clear-all",
							Value: fmt.Sprintf(
								"Choose one or many channel(s) and purge all messages.\n Usage: %s",
								"/clear-all",
							),
						},
						{
							Name: "generate",
							Value: fmt.Sprintf(
								"Recommends a new profile picture, pass in a tag like 'high_waist_skirt' and receive a new picture.\n Usage: %s",
								"/generate <tag> <tag> <nsfw>",
							),
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    fmt.Sprintf("Requested by %s", ic.Member.User.Username),
						IconURL: base.IconURL,
					},
					Timestamp: fmt.Sprint(time.Now().Format(time.RFC3339)),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		base.ThrowSimpleInteractionError(s, ic, err.Error())
	}
}
