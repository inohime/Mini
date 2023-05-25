package generatecmd

import (
	"fmt"
	"log"
	base "main/src/ops"
	"main/src/utils"
	"main/src/utils/embed"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
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
	optMap := base.ComposeOptions(ic)

	defaultURL := "https://safebooru.donmai.us"

	if _, ok := optMap["nsfw"]; ok {
		defaultURL = "https://danbooru.donmai.us"
	}

	imgURL := fmt.Sprintf("%s/posts/random?tags=%s",
		defaultURL,
		utils.EncodeString(optMap["tag-1"].StringValue()),
	)

	if opt, ok := optMap["tag-2"]; ok {
		imgURL += fmt.Sprintf("+%s", utils.EncodeString(opt.StringValue()))
	}

	log.Println(
		base.PrintCyan("%s#%s %s %s",
			ic.Member.User.Username,
			ic.Member.User.Discriminator,
			color.HiWhiteString("requested url:"),
			imgURL,
		),
	)

	img := acquireImgData(imgURL)
	if img == nil {
		base.ThrowInteractionError(
			s, ic,
			"Error finding image!",
			"Make sure the tag(s) exist(s) and is formatted properly!\nEx.) long hair -> long_hair ✅",
		)
		return
	}

	gtr, _ := img.Data.Read("general")
	generalTags := utils.EvictChars(strings.Join(gtr.([]string), ", "))

	ctr, _ := img.Data.Read("characters")
	markedChars := utils.StringsToMarkup(ctr.([]string), defaultURL)
	charTags := utils.EvictChars(strings.Join(markedChars, ", "))

	atr, _ := img.Data.Read("artist")
	markedArtists := utils.StringsToMarkup(atr.([]string), defaultURL)
	artistTags := utils.EvictChars(strings.Join(markedArtists, ", "))

	image, _ := img.Data.Read("image")
	imgSrc, _ := img.Data.Read("imgsrc")

	genMsgEmbed := embed.New(false).
		SetTitle("Recommendation for "+ic.Member.User.Username).
		SetImage(image.([]string)[0]).
		SetColor(utils.RandomColor()).
		SetField("Artist", artistTags, false).
		SetField("Tags", generalTags, false).
		SetField("Character(s)", charTags, false).
		SetField("Source", imgSrc.([]string)[0], false).
		SetFooter(fmt.Sprintf("Requested by %s", ic.Member.User.Username), base.IconURL).
		SetTimestamp(fmt.Sprint(time.Now().Format(time.RFC3339))).
		Bind()

	genMsgEmbed.Use(genMsgEmbed.DeferredResponse, s, ic).With(base.ThrowSimpleInteractionError)
}

func acquireImgData(uri string) *utils.Tags {
	doc, err := utils.FetchPageNode(uri)
	if err != nil {
		log.Println(
			base.PrintRed(
				"Failed to fetch node from page: %s",
				base.PrintWhite(err.Error()),
			),
		)
		return nil
	}

	tags := utils.NewTag(doc)
	if tags == nil {
		log.Println(base.PrintRed("Failed to create new tag object"))
		return nil
	}

	var swg sync.WaitGroup

	swg.Add(5)
	go tags.FindArtistName(&swg)
	go tags.FindGeneralTags(&swg)
	go tags.FindImageUrl(&swg)
	go tags.FindImageSource(&swg)
	go tags.FindCharacters(&swg)
	swg.Wait()

	return tags
}
