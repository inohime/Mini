package base

import (
	"fmt"
	"main/src/utils/embed"
	store "main/src/utils/store"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

var (
	IconURL = "https://cdn.discordapp.com/avatars/1002274542737182871/c06dd02b3f83235e3fe33e3fea72f7ef.png?size=1024"

	PrintGreen = color.New(color.FgHiGreen).SprintfFunc()
	PrintCyan  = color.New(color.FgHiCyan).SprintfFunc()
	PrintRed   = color.New(color.FgHiRed).SprintfFunc()
	PrintWhite = color.New(color.FgHiWhite).SprintFunc()

	Store = store.New()
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

type IBaseComponent interface {
	Name() string // used for CustomID
	Execute(*discordgo.Session, *discordgo.InteractionCreate)
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
	errMsgEmbed := embed.New(true).
		SetTitle(title).
		SetDescription(desc).
		SetImage("https://pbs.twimg.com/media/FZ8WhlkXkAAug7p?format=png&name=large").
		SetField(
			".*･｡ﾟ*☆",
			fmt.Sprintf(
				"Artwork by: [%s](%s) 🎀",
				"chromuchromu",
				"https://twitter.com/chromuchromu/",
			),
			false,
		).
		SetFooter(fmt.Sprintf("Requested by %s", ic.Member.User.Username), IconURL).
		SetTimestamp(fmt.Sprint(time.Now().Format(time.RFC3339))).
		Bind()

	errMsgEmbed.Use(errMsgEmbed.DeferredResponse, s, ic).With(ThrowSimpleInteractionError)
}

func ThrowSimpleInteractionError(s *discordgo.Session, ic *discordgo.InteractionCreate, err string) {
	_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("An error occurred: %s", err),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
