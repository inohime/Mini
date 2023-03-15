package generatecmd

import (
	"fmt"
	"main/src/utils"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type GenerateCommand struct{}

func New() *GenerateCommand {
	return &GenerateCommand{}
}

func (*GenerateCommand) Name() string {
	return "generate"
}

func (*GenerateCommand) Description() string {
	return "Generates a new profile picture"
}

func (*GenerateCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	tags := acquireTagsData()

	generalTags := strings.Join(tags.Data["general"][:], ", ")

	markedChars := utils.StringsToMarkup(tags.Data["characters"], "https://danbooru.donmai.us/")
	charTags := strings.Join(markedChars[:], ", ")

	footerEmbed := discordgo.MessageEmbedFooter{
		Text:    "Requested at: " + time.Now().Local().Format(time.ANSIC),
		IconURL: "https://cdn.discordapp.com/avatars/1002274542737182871/7ffe57c0407be99b317a78a83cb7b748.png?size=1024",
	}

	imageEmbed := discordgo.MessageEmbedImage{
		URL: tags.Data["image"][0],
	}

	s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Recommendation for " + ic.Member.User.Username,
					Image: &imageEmbed,
					Color: utils.RandomColor(),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name: "Artist", // support putting multiple artists in here
							Value: fmt.Sprintf(
								"[%s](https://danbooru.donmai.us/posts?tags=%s&z=1)",
								tags.Data["artist"][0],
								tags.Data["artist"][0],
							),
						},
						{
							Name:  "Tags",
							Value: generalTags,
						},
						{
							Name:  "Character(s)",
							Value: charTags,
						},
						{
							Name:  "Source",
							Value: tags.Data["imgsrc"][0],
						},
					},
					Footer: &footerEmbed,
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})
}

func acquireTagsData() *utils.Tags {
	doc := utils.FetchPageNode("https://danbooru.donmai.us/posts/6145011")

	tags := utils.Tags{
		Data: make(map[string][]string, 6),
		Node: doc,
	}

	tags.Wg.Add(6)
	go tags.FindArtistName()
	go tags.FindGeneralTags()
	go tags.FindImageUrl()
	go tags.FindImageSource()
	go tags.FindCharacters()
	go tags.FindCopyright()
	tags.Wg.Wait()

	return &tags
}
