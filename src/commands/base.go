package basecmd

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	IconURL = "https://cdn.discordapp.com/avatars/1002274542737182871/c06dd02b3f83235e3fe33e3fea72f7ef.png?size=1024"
)

type IBaseCommand interface {
	Name() string
	Description() string
	Execute(*discordgo.Session, *discordgo.InteractionCreate)
}

type IBaseCommandEx interface {
	IBaseCommand
	Options() []*discordgo.ApplicationCommandOption
}

func ComposeOptions(ic *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	iopts := ic.ApplicationCommandData().Options
	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(iopts))

	for _, opt := range iopts {
		optMap[opt.Name] = opt
	}

	return optMap
}

func ThrowInteractionError(s *discordgo.Session, ic *discordgo.InteractionCreate, title, desc string) {
	_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Description: desc,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name: ".*･｡ﾟ*☆",
							Value: fmt.Sprintf("Artwork by: [%s](%s) 🎀",
								"chromuchromu",
								"https://twitter.com/chromuchromu/",
							),
						},
					},
					Image: &discordgo.MessageEmbedImage{
						URL: "https://pbs.twimg.com/media/FZ8WhlkXkAAug7p?format=png&name=large",
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    fmt.Sprintf("Requested by %s", ic.Member.User.Username),
						IconURL: IconURL,
					},
					Timestamp: fmt.Sprint(time.Now().Format(time.RFC3339)),
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})
}
