package clearallcmd

import (
	"fmt"
	"log"
	base "main/src/ops"
	"main/src/utils"
	"time"

	"github.com/bwmarrin/discordgo"
)

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

func (c *ClearAllCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	userPerms := ic.Member.Permissions
	if userPerms&discordgo.PermissionAdministrator == 0 {
		base.ThrowInteractionError(
			s, ic,
			"Missing Permissions!",
			"Make sure you have elevated permissions to use this command",
		)
		return
	}

	guildChannels, err := s.GuildChannels(ic.GuildID)
	if err != nil {
		log.Println(
			base.PrintRed(
				"An error has occurred: %s",
				base.PrintWhite(err.Error()),
			),
		)
	}

	channels := make(map[string]string)
	minMenuOpt := 1

	for _, channel := range guildChannels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}
		channels[channel.Name] = channel.ID
	}

	err = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Select the channels you want to purge",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "clear-channel",
							Placeholder: "Purge these channels..",
							MinValues:   &minMenuOpt,
							MaxValues:   len(channels),
							Options:     menuOptions(channels),
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

	// wait for the menu options to be selected and event fire
	for !base.Store.ViewMenuState() {
		time.Sleep(time.Second)
	}
	base.Store.SetMenuState(false)

	// prep for the channels selected to be purged
	channelIDs, ok := base.Store.Acquire("--cccomp-channelIDs").([]string)
	if !ok {
		log.Println(base.PrintRed("Failed to acquire channelIDs"))
		return
	}

	var msgIDs []string
	var oldMsgIDs []*discordgo.Message

	// purge messages in all selected channels
	for _, id := range channelIDs {
		msgs, _ := s.ChannelMessages(id, 100, "", "", "")
		for _, msg := range msgs {
			if time.Since(msg.Timestamp).Hours() > 336 {
				oldMsgIDs = append(oldMsgIDs, msg)
			} else {
				msgIDs = append(msgIDs, msg.ID)
			}
		}

		err = s.ChannelMessagesBulkDelete(id, msgIDs)
		if err != nil {
			log.Println("Failed to delete all messages in this channel:", err)
		}

		for _, msg := range oldMsgIDs {
			err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
			if err != nil {
				log.Println("Failed to delete message:", err)
				break
			}
		}
	}

	// clear out the previous menu interaction
	s.InteractionResponseDelete(ic.Interaction)

	msg, err := s.FollowupMessageCreate(ic.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Color:       utils.RandomColor(),
				Description: "Cleared! ✅",
				Footer: &discordgo.MessageEmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", ic.Member.User.Username),
					IconURL: base.IconURL,
				},
				Timestamp: fmt.Sprint(time.Now().Format(time.RFC3339)),
			},
		},
	})
	if err != nil {
		_, _ = s.FollowupMessageCreate(ic.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("An error occurred: %s", err.Error()),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	time.Sleep(time.Second * 15)

	s.FollowupMessageDelete(ic.Interaction, msg.ID)
}

func menuOptions(channels map[string]string) []discordgo.SelectMenuOption {
	opts := make([]discordgo.SelectMenuOption, 0, len(channels))
	for channelName, channelID := range channels {
		opts = append(opts, discordgo.SelectMenuOption{
			Label: channelName,
			// optionally, switch to ToLowerSpecial to handle unicode channel names
			Value: channelID,
			Emoji: discordgo.ComponentEmoji{
				Name: "🗒️",
			},
		})
	}

	return opts
}
