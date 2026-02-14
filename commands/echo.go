package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Automaticamente registamos o comando e especificamos os dados e restrições
func init() {
	AddCommand(func(cfg *Config) BotCommand {
		return BotCommand{
			Definition: &discordgo.ApplicationCommand{
				Name:        "echo",
				Description: "Repete o que é inserido",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type: discordgo.ApplicationCommandOptionString,
						Name: "texto", Description: "Texto a ser repetido",
						Required: true,
					},
				},
			},
			Handler:       EchoHandler,
			RequiredRoles: []string{cfg.RoleDev, cfg.RoleDirecao},
			Cooldown:      5 * time.Second,
		}
	})
}

// EchoHandler contém a lógica do comando
func EchoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Extraimos o texto do comando, a mensagem a ser repetida vai estar na primeira posição
	msg := i.ApplicationCommandData().Options[0].StringValue()

	// Output do comando
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	})
}
