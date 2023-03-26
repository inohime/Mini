package viewscomp

import (
	"log"
	base "main/src/ops"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ViewLeftComponent struct{}

func NewViewLeft() *ViewLeftComponent {
	return &ViewLeftComponent{}
}

func (*ViewLeftComponent) Name() string {
	return "tags-view-left"
}

func (vl *ViewLeftComponent) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	updateView(s, ic, DECREMENT, func(num, _ int) bool {
		if num < 1 {
			log.Println(
				base.PrintRed(
					"Error: under limit %s", base.PrintWhite("searchNum bounds check"),
				),
			)

			base.ThrowSimpleInteractionError(s, ic, "You are already on the first page!")
			time.AfterFunc(time.Second*5, func() {
				_ = s.InteractionResponseDelete(ic.Interaction)
			})

			return false
		}

		return true
	})
}
