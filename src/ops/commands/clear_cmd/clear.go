package clearcmd

import (
	"fmt"
	"log"
	base "main/src/ops"
	"main/src/utils"
	"main/src/utils/embed"
	"time"

	"github.com/bwmarrin/discordgo"
)

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

func (*ClearCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "channel",
			Description: "Clear all messages in a channel",
			Required:    true,
		},
	}
}

func (*ClearCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	// overkill since there is only one option but more readable
	optMap := base.ComposeOptions(ic)
	channel := optMap["channel"].ChannelValue(s).ID

	userPerms := ic.Member.Permissions
	if userPerms&discordgo.PermissionAdministrator == 0 {
		base.ThrowInteractionError(
			s, ic,
			"Missing Permissions!",
			"Make sure you have elevated permissions to use this command",
		)
		return
	}

	var msgIDs []string
	var oldMsgIDs []*discordgo.Message

	msgs, _ := s.ChannelMessages(channel, 100, "", "", "")
	for _, msg := range msgs {
		if time.Since(msg.Timestamp).Hours() > 336 {
			oldMsgIDs = append(oldMsgIDs, msg)
		} else {
			msgIDs = append(msgIDs, msg.ID)
		}
	}

	err := s.ChannelMessagesBulkDelete(channel, msgIDs)
	if err != nil {
		log.Println(
			base.PrintRed("Failed to delete all messages in this channel: %s", err),
		)
		return
	}

	for _, msg := range oldMsgIDs {
		err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if err != nil {
			log.Println(
				base.PrintRed("Failed to delete message: %s", err),
			)
			break
		}
	}

	clearMsgEmbed := embed.New(false).
		SetColor(utils.RandomColor()).
		SetDescription("Cleared! ✅").
		SetFooter(fmt.Sprintf("Requested by %s", ic.Member.User.Username), base.IconURL).
		SetTimestamp(fmt.Sprint(time.Now().Format(time.RFC3339))).
		Bind()

	clearMsgEmbed.Use(clearMsgEmbed.Response, s, ic).With(base.ThrowSimpleInteractionError)

	time.Sleep(time.Second * 15)

	clearMsgEmbed.Use(clearMsgEmbed.ResponseDelete, s, ic)
}
