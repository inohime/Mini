package synthetic

import (
	"encoding/json"
	"fmt"
	"log"
	basecmd "main/src/commands"
	carrotcmd "main/src/commands/carrot_cmd"
	clearallcmd "main/src/commands/clear_all_cmd"
	clearcmd "main/src/commands/clear_cmd"
	generatecmd "main/src/commands/generate_cmd"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

var (
	printRed   = color.New(color.FgHiRed).SprintfFunc()
	printWhite = color.New(color.FgHiWhite).SprintFunc()
)

type Synthetic struct {
	Name          string `json:"name"`
	Token         string `json:"token"`
	Session       *discordgo.Session
	Commands      []*discordgo.ApplicationCommand
	_intlCommands map[string]basecmd.IBaseCommand // if necessary, make a commandhandler
}

func New(filePath string) (*Synthetic, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf(printWhite(err))
	}

	var synthetic Synthetic

	err = json.Unmarshal(bytes, &synthetic)
	if err != nil {
		return nil, fmt.Errorf(printWhite(err))
	}

	synthetic._intlCommands = make(map[string]basecmd.IBaseCommand)
	synthetic.Session, err = discordgo.New("Bot " + synthetic.Token)
	if err != nil {
		return nil, fmt.Errorf(printWhite(err))
	}

	return &synthetic, nil
}

func (synth *Synthetic) AddCommand(cmd basecmd.IBaseCommand) {
	synth._intlCommands[cmd.Name()] = cmd
	synth.Commands = append(synth.Commands, &discordgo.ApplicationCommand{
		Name:        cmd.Name(),
		Description: cmd.Description(),
	})
}

func (synth *Synthetic) BindCommands() {
	for i, cmd := range synth.Commands {
		log.Printf("Binding command %v...", cmd.Name)
		handle, err := synth.Session.ApplicationCommandCreate(synth.Session.State.User.ID, "", cmd)
		if err != nil {
			panic(
				printRed(
					"Failed to create %v%v%v: %v",
					printRed("["), printWhite(cmd.Name), printRed("]"),
					printWhite(err),
				),
			)
		}
		// replace the elements from the 0th index and up with the handle given
		synth.Commands[i] = handle
	}
}

func (synth *Synthetic) UnbindCommands() {
	for _, cmd := range synth.Commands {
		log.Printf("Unbinding command %v...", cmd.Name)
		err := synth.Session.ApplicationCommandDelete(synth.Session.State.User.ID, "", cmd.ID)
		if err != nil {
			panic(
				printRed(
					"Failed to delete %v%v%v: %v",
					printRed("["), printWhite(cmd.Name), printRed("]"),
					printWhite(err),
				),
			)
		}
	}
}

func (synth *Synthetic) SetupHandlers() {
	log.Println("Setting up handlers..")

	synth.Session.AddHandler(func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
		if cmd, ok := synth._intlCommands[ic.ApplicationCommandData().Name]; ok {
			cmd.Execute(s, ic)
		}
	})
	synth.Session.AddHandler(func(s *discordgo.Session, ready *discordgo.Ready) {
		log.Printf(
			"%v#%v %v",
			synth.Session.State.User.Username,
			synth.Session.State.User.Discriminator,
			color.HiGreenString("logged in"),
		)
		synth.Session.UpdateListeningStatus("your requests~!")
	})
}

// make a test for searching with html elementnodes

func Boot() {
	// create a new synthetic instance
	synthetic, err := New("src/synthetic.json")
	if err != nil {
		panic(
			printRed("Failed to create bot: %v", printWhite(err)),
		)
	}

	// add all of the commands for the bot
	synthetic.AddCommand(carrotcmd.New())
	synthetic.AddCommand(clearcmd.New())
	synthetic.AddCommand(clearallcmd.New())
	synthetic.AddCommand(generatecmd.New())

	// create and setup handlers
	synthetic.SetupHandlers()

	// open a new socket connection to discord
	err = synthetic.Session.Open()
	if err != nil {
		panic(
			printRed("Failed to open websocket: %v", printWhite(err)),
		)
	}
	defer synthetic.Session.Close()

	// create the bot's application commands
	synthetic.BindCommands()
	defer synthetic.UnbindCommands()

	log.Printf("%v is awake! Press Ctrl+C to sleep", synthetic.Name)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sig

	log.Printf("%v is going to sleep~", synthetic.Name)
}
