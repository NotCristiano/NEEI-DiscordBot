package commands

import (
	"NEEI-DiscordBot/internal/config"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	AddCommand(func(cfg *config.Config) BotCommand {
		return BotCommand{

			Definition: &discordgo.ApplicationCommand{
				Name:        "kick",
				Description: "Expulsa um user especifico",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type: discordgo.ApplicationCommandOptionUser,
						Name: "membro", Description: "Membro a ser expulso",
						Required: true,
					},
				},
			},

			Handler:       kickHandler,
			RequiredRoles: []string{cfg.RoleDev, cfg.RoleDirecao},
			Cooldown:      5 * time.Second,
			Ephemeral:     true,
		}
	})
}

// kickHandler contém a lógica do comando
func kickHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Detetamos o user a ser expulso nas options
	user := i.ApplicationCommandData().Options[0].UserValue(s)

	// Verifica se o ID de user a sere expulso é o mesmo do User que invocou a interação
	if user.ID == i.Member.User.ID {
		EditDeferredResponse(s, i, "Não te podes expulsar a ti próprio.")
		return
	}

	// Execução comando
	err := s.GuildMemberDelete(i.GuildID, user.ID)
	if err != nil {
		EditDeferredResponse(s, i, "ERRO: Falha ao expulsar membro.")
		return
	}

	EditDeferredResponse(s, i, "Membro expulso com sucesso.")
}
