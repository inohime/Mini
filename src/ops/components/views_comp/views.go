package viewscomp

import (
	"fmt"
	base "main/src/ops"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Views struct {
	IC         *discordgo.InteractionCreate
	OriginalIC *discordgo.InteractionCreate
	SearchNum  int
	TagsList   []string
	MaxLen     int
}

// Creates a new view for the Views Component (view the list of tags from left to right)
func NewViews(ic *discordgo.InteractionCreate) *Views {
	return &Views{
		IC: ic,
	}
}

// Waits until the key is found in the store before assigning it a value
//
// # Use only if there is a key in the store
//
// Assigns a value to item
func findAny[T any](key string, item *T) {
	for base.Store.Acquire(key) == nil {
		time.Sleep(time.Second)
	}

	*item = base.Store.Acquire(key).(T)
}

// Attempts to get the original message interaction
//
// Takes in the key: --tcomp-original-ic-guildID-userID
//
// Assigns the original message interaction unique to the user to Views.OriginalIC
func (v *Views) findTagsInteraction(swg *sync.WaitGroup) {
	defer swg.Done()

	originalIC := fmt.Sprintf("--tcomp-original-ic-%s-%s", v.IC.GuildID, v.IC.Member.User.ID)
	findAny(originalIC, &v.OriginalIC)
}

// Attempts to get the current tag view search number
//
// Takes in the key: --tcomp-search-num-guildID-userID
//
// The guildID and userID is to make each user's search number unique
// Assigns the unique search number to Views.SearchNum
func (v *Views) findSearchNumber(swg *sync.WaitGroup) {
	defer swg.Done()

	searchNum := fmt.Sprintf("--tcomp-search-num-%s-%s", v.IC.GuildID, v.IC.Member.User.ID)
	findAny(searchNum, &v.SearchNum)
}

// Attempts to get the list of tags for our view
//
// Takes in the key: --tcomp-tags-list-guildID-userID
//
// Assigns the unique list of tags requested by the user to Views.TagsList
func (v *Views) findTagsList(swg *sync.WaitGroup) {
	defer swg.Done()

	tagsList := fmt.Sprintf("--tcomp-tags-list-%v-%v", v.IC.GuildID, v.IC.Member.User.ID)
	findAny(tagsList, &v.TagsList)
}

// Attempts to get the maximum length of a string for formatting the list
//
// Takes in the key: --tcomp-tag-max-length-guildID-userID
//
// Assigns the string max length to Views.MaxLen
func (v *Views) findStringMaxLength(swg *sync.WaitGroup) {
	defer swg.Done()

	maxLen := fmt.Sprintf("--tcomp-tag-max-length-%v-%v", v.IC.GuildID, v.IC.Member.User.ID)
	findAny(maxLen, &v.MaxLen)
}
