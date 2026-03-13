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
				Name:        "ban",
				Description: "Banimento de um user especifico",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type: discordgo.ApplicationCommandOptionUser,
						Name: "membro", Description: "Membro a ser banido",
						Required: true,
					},
					{
						Type: discordgo.ApplicationCommandOptionInteger,
						Name: "dias", Description: "Dias que o membro vai ser banido",
						Required: true,
					},
					{
						Type: discordgo.ApplicationCommandOptionString,
						Name: "razao", Description: "Razao pela qual membro foi banido",
						Required: true,
					},
				},
			},

			Handler:       banHandler,
			RequiredRoles: []string{cfg.RoleDev, cfg.RoleDirecao},
			Cooldown:      5 * time.Second,
		}
	})
}

// kickHandler contém a lógica do comando
func banHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Detetamos o user a ser expulso e os tempo de ban nas options
	user := i.ApplicationCommandData().Options[0].UserValue(s)
	days := int(i.ApplicationCommandData().Options[1].IntValue())
	reason := i.ApplicationCommandData().Options[2].StringValue()

	// Verifica se o ID de user a ser banido é o mesmo do User que invocou a interação
	if user.ID == i.Member.User.ID {
		SendEphemeral(s, i, "Não te podes banir a ti próprio.")
		return
	}

	// Verifica tempo de banimento
	if days <= 0 {
		SendEphemeral(s, i, "O tempo de ban tem que ser superior a 0 dias")
		return
	}

	if len(reason) <= 0 {
		SendEphemeral(s, i, "Não podes banir um membro sem razão")
		return
	}

	// Execução comando
	err := s.GuildBanCreateWithReason(i.GuildID, user.ID, reason, days)
	if err != nil {
		SendEphemeral(s, i, "ERRO: Falha ao banir membro.")
		return
	}

	SendEphemeral(s, i, "Membro foi banido com sucesso.")
}
