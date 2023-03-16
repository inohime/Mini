package clearallcmd

import "github.com/bwmarrin/discordgo"

type ClearAllCommand struct{}

func New() *ClearAllCommand {
	return &ClearAllCommand{}
}

func (*ClearAllCommand) Name() string {
	return "clear-all"
}

func (*ClearAllCommand) Description() string {
	return "Purges all messages in every channel"
}

func (*ClearAllCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Carrot!",
		},
	})
}
