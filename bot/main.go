package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	db "github.com/Oluwatunmise-olat/Horus-Clone/db"
	"github.com/bwmarrin/discordgo"
)

var (
	slashCommands = []*discordgo.ApplicationCommand{
		{
			Name: "welcome",
			Description: "Welcome message",
		},
		{
			Name: "logs",
			Description: "View registered server route logs",
			Options: []*discordgo.ApplicationCommandOption{
			{
				Name: "method",
				Description: "Http request method",
				Required: false,
				Type: discordgo.ApplicationCommandOptionInteger,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name: "GET",
						Value: 1,
					},
					{	
						Name: "POST",
						Value: 2,
						},
					},
				},
			},
		},
	}

	slashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"welcome": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "welcome, to horus clone ⚡︎",
				},
			})
		},
		"logs": func(session *discordgo.Session, i *discordgo.InteractionCreate) {
			choice := i.ApplicationCommandData().Options
			var optionValue int64

			if (len(choice) > 0) {
				optionValue = choice[0].IntValue()
			}

			method := ""

			if optionValue == 1 {
				method = "get"
			}else if optionValue == 2{
				method = "post"
			}else {
			}

			dbLimit := 1 // Limit to one record because of timeout
			data := Conf.Db.GetLogs(&dbLimit, method)

			resp := ""

			for _, record := range data {
				stringify, _ := json.Marshal(record)
				resp += string(stringify) + ","
			}

			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: resp,
				},
			})
		},
	}

	Conf *DiscordConfig
)

type DiscordConfig struct {
	AppId string
	Db *db.Store
	Token string
	GuildId string

	session *discordgo.Session
	storedCommands []*discordgo.ApplicationCommand
}

func Init(d *DiscordConfig){
	// Respond to commands
	d.session.AddHandler(func (s *discordgo.Session, i *discordgo.InteractionCreate){
		if commandHandler, ok := slashCommandHandlers[i.ApplicationCommandData().Name]; ok {
			commandHandler(s, i)
		}
	})

	// Register commands
	storedCommands := make([]*discordgo.ApplicationCommand, len(slashCommands))

	for index, value := range slashCommands {
		log.Printf("Creating %q command", value.Name)
		cmd, err := d.session.ApplicationCommandCreate(Conf.AppId, Conf.GuildId, value);

		if err != nil {
			log.Fatalf("Could not create %q command, Error: %v", value.Name, err)
		}
		
		storedCommands[index] = cmd
	}

	d.storedCommands = storedCommands

	// Open websocket session to discord
	err := d.session.Open()
	if err != nil {
		log.Fatalf("Could not open ws session: %v", err)
	}
	defer d.session.Close()
}

func cleanUp(storedCommands []*discordgo.ApplicationCommand, s *discordgo.Session) {
	// Delete all created bot commands
	for _, value := range storedCommands {
		s.ApplicationCommandDelete(Conf.AppId, Conf.GuildId, value.ID)
	}
}

func Connect(dc *DiscordConfig) error {
	session, err := discordgo.New("Bot " + dc.Token)

	if err != nil {
		return fmt.Errorf("Could not connect to discord. ¦")
	}

	dc.session = session	
	Conf = dc
	Init(dc)

	return nil
}


func ListenForInterrupt(){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cleanUp(Conf.storedCommands, Conf.session)
	os.Exit(0)
}