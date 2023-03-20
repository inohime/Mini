package basecmd

import "github.com/bwmarrin/discordgo"

type IBaseCommand interface {
	Name() string
	Description() string
	Execute(*discordgo.Session, *discordgo.InteractionCreate)
}

type IBaseCommandEx interface {
	IBaseCommand
	Options() []*discordgo.ApplicationCommandOption
}
