package generatecmd

import (
	"fmt"
	"log"
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

func (*GenerateCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tag-1",
			Description: "A label used to fetch a specific image: Ex.) wide hips -> wide_hips ✅",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tag-2",
			Description: "A label used to fetch a specific image: Ex.) wide hips -> wide_hips ✅",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "nsfw",
			Description: "Allows searching for very explicit content",
			Required:    false,
		},
	}
}

func (*GenerateCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	iopts := ic.ApplicationCommandData().Options
	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(iopts))

	for _, opt := range iopts {
		optMap[opt.Name] = opt
	}

	defaultURL := "https://safebooru.donmai.us"

	if _, ok := optMap["nsfw"]; ok {
		defaultURL = "https://danbooru.donmai.us"
	}

	iconURL := "https://cdn.discordapp.com/avatars/1002274542737182871/c06dd02b3f83235e3fe33e3fea72f7ef.png?size=1024"
	imgURL := fmt.Sprintf("%s/posts/random?tags=%s", defaultURL, utils.EncodeString(optMap["tag-1"].StringValue()))

	if opt, ok := optMap["tag-2"]; ok {
		imgURL += fmt.Sprintf("+%s", utils.EncodeString(opt.StringValue()))
	}

	log.Println("the url requested:", imgURL)

	img := acquireImgData(imgURL)
	if img == nil {
		_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Error finding image!",
						Description: "Make sure the tag(s) exist(s) and is formatted properly!\nEx.) long hair -> long_hair ✅",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  ".*･｡ﾟ*☆",
								Value: fmt.Sprintf("Artwork by: [%s](%s) 🎀", "chromuchromu", "https://twitter.com/chromuchromu/"),
							},
						},
						Image: &discordgo.MessageEmbedImage{
							URL: "https://pbs.twimg.com/media/FZ8WhlkXkAAug7p?format=png&name=large",
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text:    fmt.Sprintf("<t:%v>", time.Now().Unix()),
							IconURL: iconURL,
						},
					},
				},
				AllowedMentions: &discordgo.MessageAllowedMentions{},
			},
		})

		return
	}

	generalTags := strings.Join(img.Data["general"][:], ", ")

	markedChars := utils.StringsToMarkup(img.Data["characters"], defaultURL)
	charTags := strings.Join(markedChars[:], ", ")

	markedArtists := utils.StringsToMarkup(img.Data["artist"], defaultURL)
	artistTags := strings.Join(markedArtists[:], ", ")

	err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Recommendation for " + ic.Member.User.Username,
					Image: &discordgo.MessageEmbedImage{
						URL: img.Data["image"][0],
					},
					Color: utils.RandomColor(),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Artist",
							Value: artistTags,
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
							Value: img.Data["imgsrc"][0],
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    fmt.Sprintf("%v requested", ic.Member.User.Username),
						IconURL: iconURL,
					},
					Timestamp: fmt.Sprint(time.Now().Format(time.RFC3339)),
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})
	if err != nil {
		_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error occured: %v", err.Error()),
			},
		})
	}
}

func acquireImgData(uri string) *utils.Tags {
	doc, err := utils.FetchPageNode(uri)
	if err != nil {
		return nil
	}

	tags := utils.NewTag(doc)
	// (de/in)crement based on the number of tag functions
	numTasks := 5

	tags.Sync.Add(numTasks)
	go tags.FindArtistName()
	go tags.FindGeneralTags()
	go tags.FindImageUrl()
	go tags.FindImageSource()
	go tags.FindCharacters()
	tags.Sync.Wait()

	return tags
}
