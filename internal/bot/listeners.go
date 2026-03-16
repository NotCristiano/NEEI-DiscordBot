package bot

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func messageCreateHandler(cfg *config.Config) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		// Ignora mensagens do próprio bot
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Ignora mensagens de bots
		if m.Author.Bot {
			return
		}

		// Se for um DM, ignora
		if m.GuildID == "" {
			return
		}

		// Ignora mensagens de users com roles especificas
		if roles := m.Member.Roles; len(roles) > 0 {
			for _, roleID := range m.Member.Roles {
				if roleID == cfg.RoleNEEI {
					return
				}
			}
		}

		// Verifica se a mensagem foi enviada no canal proibido
		if m.ChannelID == cfg.MuteChannelID {
			// Aplica o timeout ao autor da mensagem
			// ATENÇÃO: NO ENV AutoMuteDuration DEVE SER ESCRITO PARA MINUTOS
			err := commands.ApplyTimeout(s, m.GuildID, m.Author.ID, time.Duration(cfg.AutoMuteDuration)*time.Minute)
			if err != nil {
				log.Error().Err(err).Str("component", "messageCreateHandler").Msg("Erro ao aplicar timeout.")
			}
			log.Info().Str("user", m.Author.Username).Str("id", m.Author.ID).Msgf("Usuário recebeu timeout de %d minutos por enviar mensagem no canal proibido.", cfg.AutoMuteDuration)

			// Se concluiu com sucesso apagamos a mensagem
			err = s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				log.Error().Err(err).Str("component", "messageCreateHandler").Msg("Erro ao apagar mensagem.")
			}
			log.Info().Str("user", m.Author.Username).Str("id", m.Author.ID).Msg("Mensagem apagada com sucesso.")

		}
	}
}
