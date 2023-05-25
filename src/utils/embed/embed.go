package embed

import (
	"github.com/bwmarrin/discordgo"
)

// redesign this later, we need to support InteractionResponseData as a whole

type MsgEmbed struct {
	*discordgo.MessageEmbed
	Err       []error
	Ephemeral bool
}

func New(ephemeral bool) *MsgEmbed {
	msgEmbed := &MsgEmbed{
		&discordgo.MessageEmbed{},
		make([]error, 10),
		false,
	}

	if ephemeral {
		msgEmbed.Ephemeral = true
	}

	return msgEmbed
}

func (e *MsgEmbed) SetURL(url string) *MsgEmbed {
	e.URL = url
	return e
}

func (e *MsgEmbed) SetTitle(title string) *MsgEmbed {
	e.Title = title
	return e
}

func (e *MsgEmbed) SetDescription(desc string) *MsgEmbed {
	e.Description = desc
	return e
}

func (e *MsgEmbed) SetTimestamp(timestamp string) *MsgEmbed {
	e.Timestamp = timestamp
	return e
}

func (e *MsgEmbed) SetColor(col int) *MsgEmbed {
	e.Color = col
	return e
}

func (e *MsgEmbed) SetFooter(text, iconURL string) *MsgEmbed {
	e.Footer = &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return e
}

func (e *MsgEmbed) SetImage(url string) *MsgEmbed {
	e.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
	return e
}

func (e *MsgEmbed) SetThumbnail(url string, width, height int) *MsgEmbed {
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:    url,
		Width:  width,
		Height: height,
	}
	return e
}

func (e *MsgEmbed) SetAuthor(url, name, iconURL string) *MsgEmbed {
	e.Author = &discordgo.MessageEmbedAuthor{
		URL:     url,
		Name:    name,
		IconURL: iconURL,
	}
	return e
}

func (e *MsgEmbed) SetField(name, value string, inline bool) *MsgEmbed {
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return e
}

func (e *MsgEmbed) SetFieldEx(fields ...*discordgo.MessageEmbedField) *MsgEmbed {
	e.Fields = fields
	return e
}

func (e *MsgEmbed) Use(
	fn func(*discordgo.Session, *discordgo.InteractionCreate),
	s *discordgo.Session,
	ic *discordgo.InteractionCreate,
) *MsgEmbed {
	fn(s, ic)
	return e
}

// for now, this is just for command error handling
func (e *MsgEmbed) With(fn func(*discordgo.Session, *discordgo.InteractionCreate, string)) *MsgEmbed {
	_ = func(s *discordgo.Session, ic *discordgo.InteractionCreate, err string) {
		if len(e.Err) == 0 {
			return
		}
		fn(s, ic, err)
		// clear the list of errors
		e.Err = e.Err[:0]
	}
	return e
}

func (e *MsgEmbed) isEphemeral() discordgo.MessageFlags {
	if e.Ephemeral {
		return discordgo.MessageFlagsEphemeral
	}
	return 0
}

func (e *MsgEmbed) DeferredResponse(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: e.isEphemeral(),
		},
	})
	if err != nil {
		e.Err = append(e.Err, err)
	}

	_, err = s.FollowupMessageCreate(ic.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			e.MessageEmbed,
		},
		Flags: e.isEphemeral(),
	})
	if err != nil {
		e.Err = append(e.Err, err)
	}
}

func (e *MsgEmbed) Response(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				e.MessageEmbed,
			},
			Flags: e.isEphemeral(),
		},
	})
	if err != nil {
		e.Err = append(e.Err, err)
	}
}

func (e *MsgEmbed) ResponseDelete(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	err := s.InteractionResponseDelete(ic.Interaction)
	if err != nil {
		e.Err = append(e.Err, err)
	}
}

// Does nothing other than make the builder feel complete
func (e *MsgEmbed) Bind() *MsgEmbed {
	return e
}

func (e *MsgEmbed) Unwrap() *discordgo.MessageEmbed {
	return e.MessageEmbed
}
