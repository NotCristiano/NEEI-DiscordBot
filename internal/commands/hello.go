package commands

import (
	"NEEI-DiscordBot/internal/config"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Automaticamente registamos o comando e especificamos os dados e restrições
func init() {
	AddCommand(func(cfg *config.Config) BotCommand {
		return BotCommand{

			Definition: &discordgo.ApplicationCommand{
				Name:        "hello",
				Description: "Retorna 'Hello World!'",
				Options:     []*discordgo.ApplicationCommandOption{},
			},
			Handler:       helloHandler,
			RequiredRoles: nil,
			Cooldown:      0 * time.Second,
		}
	})
}

// helloHandler contém a lógica do comando
func helloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Hello World!"}})
}
