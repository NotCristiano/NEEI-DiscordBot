package bot

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func createTicket(s *discordgo.Session, i *discordgo.InteractionCreate, description string, categoryID string, cfg *config.Config) {
	userID := i.Member.User.ID
	guildID := i.GuildID
	userName := i.Member.User.Username

	exists, err := ticketExists(s, guildID, categoryID, userName)
	if err != nil {
		log.Error().Err(err).Str("component", "createTicket").Msg("Erro ao verificar existência de ticket.")
		commands.SendEphemeral(s, i, "Erro interno do bot, tente novamente mais tarde.")
		return
	}

	if exists {
		log.Info().Str("user", userName).Msg("Utilizador tentou criar um ticket, mas já existe um aberto.")
		commands.SendEphemeral(s, i, "Já tens um ticket aberto. Por favor fecha o ticket atual antes de abrir um novo.")
		return
	}

	channel, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     "ticket-" + userName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				// everyone
				ID:   guildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
			{
				// The user who opened the ticket
				ID:    userID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
			},
			{
				// TODO: Decidir que roles vão ter acesso a estes canais, provavelmente APE, TEC, DIR
				ID:    cfg.RoleDev,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
			},
		},
	})

	if err != nil {
		log.Error().Err(err).Str("component", "createTicket").Msg("Erro ao criar canal de ticket.")
		commands.SendEphemeral(s, i, "Ocorreu um erro ao criar o ticket. Por favor tenta novamente mais tarde.")
	}

	commands.SendEphemeral(s, i, "Ticket criado com sucesso! Dirige-te ao canal "+channel.Mention()+".")
	initTicketMessage(s, description, channel.ID)
	log.Info().Str("channel", channel.ID).Msg("Canal de ticket criado com sucesso!")
}

func initTicketMessage(s *discordgo.Session, description string, channelID string) {
	// Definimos o embed
	embed := &discordgo.MessageEmbed{
		Title:       "Ticket",
		Description: "**Descrição do ticket:**\n" + description + "\n\n" + "Se já não precisas de ajuda ou se o teu problema foi resolvido, por favor fecha o ticket clicando no botão abaixo.",
		Color:       0x970302,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Núcleo de Estudantes de Engenharia Informática",
			IconURL: "https://example.com/logo.png",
		},
	}

	// Definimos o botão
	button := discordgo.Button{
		Label:    "Fechar Ticket",
		Style:    discordgo.PrimaryButton,
		CustomID: "ticketCloseButton",
		Emoji: &discordgo.ComponentEmoji{
			Name: "✅",
		},
	}

	// Agora juntamos embed e button para enviar
	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{button},
			},
		},
	})

}

func ticketExists(s *discordgo.Session, guildID, categoryID string, userName string) (bool, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return false, err
	}

	ticketName := "ticket-" + userName
	for _, channel := range channels {
		if channel.Name == ticketName && channel.ParentID == categoryID {
			return true, nil
		}
	}
	return false, nil
}
