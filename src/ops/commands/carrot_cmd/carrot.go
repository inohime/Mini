package carrotcmd

import "github.com/bwmarrin/discordgo"

type CarrotCommand struct{}

func New() *CarrotCommand {
	return &CarrotCommand{}
}

func (*CarrotCommand) Name() string {
	return "carrot"
}

func (*CarrotCommand) Description() string {
	return "Says carrot in the current channel!"
}

func (*CarrotCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Carrot!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
