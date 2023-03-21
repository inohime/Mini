package clearchannelcomp

import (
	"log"
	base "main/src/ops"

	"github.com/bwmarrin/discordgo"
)

type ClearChannelComponent struct{}

func New() *ClearChannelComponent {
	return &ClearChannelComponent{}
}

func (*ClearChannelComponent) Name() string {
	return "clear-channel"
}

func (c *ClearChannelComponent) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	channels := ic.MessageComponentData().Values
	for _, c := range channels {
		log.Println("Channel selected for purging:", c)
	}
	base.Store.Bundle(channels, "--cccomp-channelIDs")
}
