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
			Handler:       echoHandler,
			RequiredRoles: []string{cfg.RoleDev, cfg.RoleDirecao, cfg.RolePresiDeptec, cfg.RoleVicePresiDeptec},
			Cooldown:      5 * time.Second,
			Ephemeral:     true,
		}
	})
}

// echoHandler contém a lógica do comando
func echoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Guardamos a mensagem para o bot enviar
	msg := i.ApplicationCommandData().Options[0].StringValue()
	if sendEcho(msg, s, i) != nil {
		EditDeferredResponse(s, i, "ERRO: Falha ao enviar mensagem.")
		return
	}
	EditDeferredResponse(s, i, "Mensagem enviada com sucesso.")

}

func sendEcho(msg string, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	_, err := s.ChannelMessageSend(i.ChannelID, msg)
	return err
}
