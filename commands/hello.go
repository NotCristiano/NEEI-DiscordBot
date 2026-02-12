package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Automaticamente registamos o comando e especificamos os dados e restrições
func init() {
	RegisterCommand(BotCommand{
		Definition: &discordgo.ApplicationCommand{
			Name:        "hello",
			Description: "Retorna 'Hello World!'",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
		Handler:       HelloHandler,
		RequiredRoles: nil,
		Cooldown:      5 * time.Second,
	})
}

// HelloHandler contém a lógica do comando
func HelloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Hello World!"}})
}
