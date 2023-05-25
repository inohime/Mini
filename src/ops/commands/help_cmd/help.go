package helpcmd

import (
	"fmt"
	base "main/src/ops"
	"main/src/utils"
	"main/src/utils/embed"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.MessageEmbedField{}
)

type HelpCommand struct{}

func New() *HelpCommand {
	return &HelpCommand{}
}

func (*HelpCommand) Name() string {
	return "help"
}

func (*HelpCommand) Description() string {
	return "Lists all commands available"
}

func init() {
	addCommandToList("carrot", fmt.Sprintf(
		"Creates a blob with the word 'carrot' in it.\n Usage: %s",
		"/carrot",
	))
	addCommandToList("clear", fmt.Sprintf(
		"Give it a channel and it'll purge all messages.\n Usage: %s",
		"/clear <channel>",
	))
	addCommandToList("clear-all", fmt.Sprintf(
		"Choose one or many channel(s) and purge all messages.\n Usage: %s",
		"/clear-all",
	))
	addCommandToList("generate", fmt.Sprintf(
		"Recommends a new profile picture, pass in a tag like 'high_waist_skirt' and receive a new picture.\n Usage: %s",
		"/generate <tag> <tag> <nsfw>",
	))
	addCommandToList("tags", fmt.Sprintf(
		"View the most popular tags or see what's available.\n Usage: %s",
		"/tags <char>",
	))
}

func (c *HelpCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	helpMsgEmbed := embed.New(true).
		SetTitle("Table of Commands").
		SetColor(utils.RandomColor()).
		SetFieldEx(commands...).
		SetFooter(fmt.Sprintf("Requested by %s", ic.Member.User.Username), base.IconURL).
		SetTimestamp(fmt.Sprint(time.Now().Format(time.RFC3339))).
		Bind()

	helpMsgEmbed.Use(helpMsgEmbed.Response, s, ic).With(base.ThrowSimpleInteractionError)
}

func addCommandToList(name, usage string) {
	commands = append(commands, &discordgo.MessageEmbedField{
		Name:  name,
		Value: usage,
	})
}
