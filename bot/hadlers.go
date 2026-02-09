package bot

import (
	"NEEI-DiscordBot/commands"

	"github.com/bwmarrin/discordgo"
)

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Se não for um comando, ignoramos
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	cmdName := i.ApplicationCommandData().Name
	if cmd, ok := commands.CommandMap[cmdName]; ok {

		// Verificamos se o comando quer uma role
		if len(cmd.RequiredRoles) > 0 {
			authorized := false

			for _, userRole := range i.Member.Roles {
				for _, requiredRole := range cmd.RequiredRoles {
					if userRole == requiredRole {
						authorized = true
						break
					}
				}
				if authorized {
					break
				}
			}

			if !authorized {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Acesso restrito aos cargos autorizados.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}
		}
		cmd.Handler(s, i)
	}
}
