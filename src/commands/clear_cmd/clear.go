package clearcmd

import "github.com/bwmarrin/discordgo"

type ClearCommand struct{}

func New() *ClearCommand {
	return &ClearCommand{}
}

func (*ClearCommand) Name() string {
	return "clear"
}

func (*ClearCommand) Description() string {
	return "Purges all messages in a channel"
}

func (*ClearCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Carrot!",
		},
	})
}
