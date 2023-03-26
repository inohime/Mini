package viewscomp

import (
	"fmt"
	"log"
	base "main/src/ops"
	"main/src/utils"
	"sort"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	// Initial viewable page number before it's updated
	pageNum = discordgo.Button{
		CustomID: "tags-page-num",
		Emoji: discordgo.ComponentEmoji{
			Name: "1️⃣",
		},
		Style: discordgo.SecondaryButton,
	}

	// List of keycap emojis (from 0-10) for pageNum
	pageMap = make(map[int]string)
)

const (
	INCREMENT = iota
	DECREMENT
)

func init() {
	for i := 1; i < 10; i++ {
		pageMap[i] = fmt.Sprintf("%d\ufe0f\u20e3", i)
	}

	pageMap[10] = "\U0001f51f"
}

// optional: searchNum in circular doubly linked list to loop
func updateView(s *discordgo.Session, ic *discordgo.InteractionCreate, op int, opFn func(int, int) bool) {
	views := acquireViewBundles(ic)
	if views == nil {
		base.ThrowInteractionError(
			s, ic,
			"Error finding view tags!",
			"Make sure the interaction is valid!",
		)
		return
	}

	tagsMap := make(map[int][]string)
	for i, x := range views.TagsList {
		key := i / 10
		tagsMap[key] = append(tagsMap[key], x)
	}

	if op == INCREMENT {
		views.SearchNum++
	} else if op == DECREMENT {
		views.SearchNum--
	}

	if !opFn(views.SearchNum, len(tagsMap)) {
		return
	}

	sort.Slice(views.TagsList[:], func(i, j int) bool {
		return views.TagsList[i] < views.TagsList[j]
	})

	updateNum(views.SearchNum)

	tagsTable := utils.MakeTagsTable(tagsMap[views.SearchNum-1], &views.MaxLen)
	title := "Table of Tags"
	title = fmt.Sprintf("%[2]*[1]s\n", title, (views.MaxLen+len(title))/2)
	content := fmt.Sprintf("```%s\n%s```", title, tagsTable)

	_, err := s.InteractionResponseEdit(views.OriginalIC.Interaction, &discordgo.WebhookEdit{
		Content: &content,
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: useNewButtons(),
			},
		},
	})
	if err != nil {
		log.Println(base.PrintRed("Failed to edit message: %s", err.Error()))
	}

	_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	updateBundles(views)
}

func acquireViewBundles(ic *discordgo.InteractionCreate) *Views {
	views := NewViews(ic)
	if views == nil {
		log.Println(base.PrintRed("Failed to create new views object"))
		return nil
	}

	var swg sync.WaitGroup

	swg.Add(4)
	go views.findTagsInteraction(&swg)
	go views.findSearchNumber(&swg)
	go views.findTagsList(&swg)
	go views.findStringMaxLength(&swg)
	swg.Wait()

	return views
}

func updateBundles(v *Views) {
	base.Store.Bundle(
		v.SearchNum,
		fmt.Sprintf("--tcomp-search-num-%s-%s", v.IC.GuildID, v.IC.Member.User.ID),
	)
	base.Store.Bundle(
		v.MaxLen,
		fmt.Sprintf("--tcomp-tag-max-length-%v-%v", v.IC.GuildID, v.IC.Member.User.ID),
	)
	base.Store.Bundle(
		v.TagsList,
		fmt.Sprintf("--tcomp-tags-list-%v-%v", v.IC.GuildID, v.IC.Member.User.ID),
	)
}

func useNewButtons() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
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
		pageNum,
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
	}
}

func updateNum(x int) {
	pageNum = discordgo.Button{
		CustomID: "tags-page-num",
		Emoji: discordgo.ComponentEmoji{
			Name: pageMap[x],
		},
		Style: discordgo.SecondaryButton,
	}
}
