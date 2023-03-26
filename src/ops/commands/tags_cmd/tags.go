package tagscmd

import (
	"encoding/json"
	"fmt"
	"log"
	base "main/src/ops"
	"main/src/utils"
	"os"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type TagResult struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	PostCount    int      `json:"post_count"`
	Category     int      `json:"category"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	IsDeprecated bool     `json:"is_deprecated"`
	Words        []string `json:"words"`
}

var (
	tagResults = make([]TagResult, 0)
)

func init() {
	bytes, err := os.ReadFile("set_of_tags.json")
	if err != nil {
		log.Panicln(
			base.PrintRed(
				"Failed to read file: %s", base.PrintWhite(err.Error()),
			),
		)
	}

	err = json.Unmarshal(bytes, &tagResults)
	if err != nil {
		log.Panicln(
			base.PrintRed(
				"Failed to unmarshal json: %v",
				base.PrintWhite(err.Error()),
			),
		)
	}

	sort.Slice(tagResults[:], func(i, j int) bool {
		return tagResults[i].Name[:1] < tagResults[j].Name[:1]
	})
}

type TagsCommand struct{}

func New() *TagsCommand {
	return &TagsCommand{}
}

func (*TagsCommand) Name() string {
	return "tags"
}

func (*TagsCommand) Description() string {
	return "Lists all of the possible tags"
}

func (*TagsCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "letter",
			Description: "choose an alphanumeric character from A-Z 0-9",
			Required:    true,
		},
	}
}

func (*TagsCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	optMap := base.ComposeOptions(ic)

	char := strings.ToLower(optMap["letter"].StringValue())
	if len(char) > 1 {
		base.ThrowSimpleInteractionError(s, ic, "Make sure your input is one character!")
		return
	}

	// filter array here
	tagsList := []string{}
	for _, x := range tagResults {
		if x.Name[:1] == char {
			tagsList = append(tagsList, x.Name)
		}
	}

	// sort all strings alphabetically
	sort.Slice(tagsList[:], func(i, j int) bool {
		return tagsList[i] < tagsList[j]
	})

	searchNum := 1
	maxLen := 0
	tagsTable := utils.MakeTagsTable(tagsList, &maxLen)
	title := "Table of Tags"
	title = fmt.Sprintf("%[2]*[1]s\n", title, (maxLen+len(title))/2)
	content := fmt.Sprintf("```%s\n%s```", title, tagsTable)

	err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: "tags-empty-first",
							Label:    "\u200c",
							Style:    discordgo.SecondaryButton,
						},
						discordgo.Button{
							CustomID: "tags-view-left",
							Label:    "<",
							Style:    discordgo.PrimaryButton,
						},
						discordgo.Button{
							CustomID: "tags-page-num",
							Emoji: discordgo.ComponentEmoji{
								Name: "1️⃣",
							},
							Style: discordgo.SecondaryButton,
						},
						discordgo.Button{
							CustomID: "tags-view-right",
							Label:    ">",
							Style:    discordgo.PrimaryButton,
						},
						discordgo.Button{
							CustomID: "tags-empty-last",
							Label:    "\u200c",
							Style:    discordgo.SecondaryButton,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		base.ThrowSimpleInteractionError(s, ic, err.Error())
	}

	base.Store.Bundle(
		searchNum,
		fmt.Sprintf("--tcomp-search-num-%v-%v", ic.GuildID, ic.Member.User.ID),
	)
	base.Store.Bundle(
		ic,
		fmt.Sprintf("--tcomp-original-ic-%v-%v", ic.GuildID, ic.Member.User.ID),
	)
	base.Store.Bundle(
		maxLen,
		fmt.Sprintf("--tcomp-tag-max-length-%v-%v", ic.GuildID, ic.Member.User.ID),
	)
	base.Store.Bundle(
		tagsList,
		fmt.Sprintf("--tcomp-tags-list-%v-%v", ic.GuildID, ic.Member.User.ID),
	)
}
