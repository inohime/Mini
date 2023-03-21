package synthetic

import (
	"encoding/json"
	"fmt"
	"log"
	base "main/src/ops"
	carrotcmd "main/src/ops/commands/carrot_cmd"
	clearallcmd "main/src/ops/commands/clear_all_cmd"
	clearcmd "main/src/ops/commands/clear_cmd"
	generatecmd "main/src/ops/commands/generate_cmd"
	clearchannelcomp "main/src/ops/components/clear_channel_comp"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

type Synthetic struct {
	Name            string `json:"name"`
	Token           string `json:"token"`
	Session         *discordgo.Session
	Commands        []*discordgo.ApplicationCommand
	_intlCommands   map[string]base.IBaseCommand
	_intlComponents map[string]base.IBaseComponent
}

func New(filePath string) (*Synthetic, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf(base.PrintWhite(err))
	}

	var synthetic Synthetic

	err = json.Unmarshal(bytes, &synthetic)
	if err != nil {
		return nil, fmt.Errorf(base.PrintWhite(err))
	}

	synthetic._intlCommands = make(map[string]base.IBaseCommand)
	synthetic._intlComponents = make(map[string]base.IBaseComponent)

	synthetic.Session, err = discordgo.New("Bot " + synthetic.Token)
	if err != nil {
		return nil, fmt.Errorf(base.PrintWhite(err))
	}

	return &synthetic, nil
}

func Boot() {
	// create a new synthetic instance
	synthetic, err := New("src/synthetic.json")
	if err != nil {
		panic(base.PrintRed("Failed to create bot: %s", base.PrintWhite(err)))
	}
	// add all of the commands for the bot
	synthetic.AddCommand(carrotcmd.New())
	synthetic.AddCommand(clearcmd.New())
	synthetic.AddCommand(clearallcmd.New())
	synthetic.AddCommand(generatecmd.New())

	// add all of the components for the bot
	synthetic.AddComponent(clearchannelcomp.New())

	// create and setup handlers
	synthetic.SetupHandlers()

	// open a new socket connection to discord
	err = synthetic.Session.Open()
	if err != nil {
		panic(base.PrintRed("Failed to open websocket: %s", base.PrintWhite(err)))
	}
	defer synthetic.Session.Close()

	// create the bot's application commands
	synthetic.BindCommands()
	defer synthetic.UnbindCommands()

	log.Printf("%s is awake! Press Ctrl+C to sleep", synthetic.Name)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sig

	log.Printf("%s is going to sleep~", synthetic.Name)
}

func (synth *Synthetic) AddCommand(cmd base.IBaseCommand) {
	synth._intlCommands[cmd.Name()] = cmd
	synth.Commands = append(synth.Commands, &discordgo.ApplicationCommand{
		Name:        cmd.Name(),
		Description: cmd.Description(),
		Options:     hasOptions(cmd),
	})
}

func (synth *Synthetic) AddComponent(comp base.IBaseComponent) {
	synth._intlComponents[comp.Name()] = comp
}

func (synth *Synthetic) BindCommands() {
	for i, cmd := range synth.Commands {
		log.Printf("Binding command %s...", cmd.Name)
		handle, err := synth.Session.ApplicationCommandCreate(synth.Session.State.User.ID, "", cmd)
		if err != nil {
			panic(
				base.PrintRed(
					"Failed to create %s%s%s: %s",
					base.PrintRed("["), base.PrintWhite(cmd.Name), base.PrintRed("]"),
					base.PrintWhite(err),
				),
			)
		}
		// replace the elements from the 0th index and up with the handle given
		synth.Commands[i] = handle
	}
}

func (synth *Synthetic) UnbindCommands() {
	for _, cmd := range synth.Commands {
		log.Printf("Unbinding command %s...", cmd.Name)
		err := synth.Session.ApplicationCommandDelete(synth.Session.State.User.ID, "", cmd.ID)
		if err != nil {
			panic(
				base.PrintRed(
					"Failed to delete %s%s%s: %s",
					base.PrintRed("["), base.PrintWhite(cmd.Name), base.PrintRed("]"),
					base.PrintWhite(err),
				),
			)
		}
	}
}

func (synth *Synthetic) SetupHandlers() {
	log.Println("Setting up handlers..")

	synth.Session.AddHandler(func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
		switch ic.Type {
		case discordgo.InteractionApplicationCommand:
			if cmd, ok := synth._intlCommands[ic.ApplicationCommandData().Name]; ok {
				cmd.Execute(s, ic)
			}
		case discordgo.InteractionMessageComponent:
			if cmd, ok := synth._intlComponents[ic.MessageComponentData().CustomID]; ok {
				cmd.Execute(s, ic)
			}

			if ic.MessageComponentData().ComponentType == discordgo.SelectMenuComponent {
				base.Store.SetMenuState(true)
			}
		}
	})
	synth.Session.AddHandler(func(s *discordgo.Session, ready *discordgo.Ready) {
		log.Println(
			base.PrintCyan("%s#%s %s",
				synth.Session.State.User.Username,
				synth.Session.State.User.Discriminator,
				color.HiGreenString("logged in"),
			),
		)
		err := synth.Session.UpdateListeningStatus("your requests~!")
		if err != nil {
			log.Println(base.PrintRed("Failed to update status:%s", base.PrintWhite(err)))
		}
	})
}

func hasOptions(cmd base.IBaseCommand) []*discordgo.ApplicationCommandOption {
	if extended, ok := cmd.(base.IBaseCommandEx); ok {
		return extended.Options()
	}
	return nil
}
