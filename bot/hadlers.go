package bot

import (
	"NEEI-DiscordBot/commands"

	"github.com/bwmarrin/discordgo"
)

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		// Procuramos o comando correspondente
		commandName := i.ApplicationCommandData().Name

		// Verificamos se o comando existe
		if handler, ok := commands.CommandMap[commandName]; ok {
			handler(s, i)
		}
	}
}
