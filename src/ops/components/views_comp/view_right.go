package viewscomp

import (
	"log"
	base "main/src/ops"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ViewRightComponent struct{}

func NewViewRight() *ViewRightComponent {
	return &ViewRightComponent{}
}

func (*ViewRightComponent) Name() string {
	return "tags-view-right"
}

func (vr *ViewRightComponent) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	updateView(s, ic, INCREMENT, func(num, len int) bool {
		if num > 10 || num > len {
			log.Println(
				base.PrintRed(
					"Error: over limit %s", base.PrintWhite("searchNum bounds check"),
				),
			)

			base.ThrowSimpleInteractionError(s, ic, "You are already on the last page!")
			time.AfterFunc(time.Second*5, func() {
				_ = s.InteractionResponseDelete(ic.Interaction)
			})

			return false
		}

		return true
	})
}
