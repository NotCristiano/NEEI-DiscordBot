package bot

import (
	"NEEI-DiscordBot/internal/commands"
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Variaveis globais para dar track no cooldown
var (
	// Formato mapa: UserID: NomeComando -> Ultima execução do comando
	userCooldowns = make(map[string]time.Time)
	cooldownMutex sync.Mutex
)

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Se não for um comando, ignoramos
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Puxamos as infos do comando
	cmdData := i.ApplicationCommandData()

	// Setup do logger e print do comando recebido
	logger := log.With().
		Str("component", "interactionHandler").
		Str("command", cmdData.Name).
		Str("user", i.Member.User.Username).
		Str("id", i.Member.User.ID).
		Str("guild", i.GuildID).
		Logger()

	logger.Debug().Msg("Comando recebido!")

	// Verificamos se o comando existe
	if cmd, ok := commands.CommandMap[cmdData.Name]; ok {

		if cmd.Cooldown > 0 {
			// Criamos uma key
			key := fmt.Sprintf("%s:%s", i.Member.User.ID, cmdData.Name)

			cooldownMutex.Lock() // Damos lock para evitar conflitos
			lastExecution, found := userCooldowns[key]

			if found {
				elapsed := time.Since(lastExecution)
				if elapsed < cmd.Cooldown {
					remaining := cmd.Cooldown - elapsed
					cooldownMutex.Unlock() // Antes de retornar temos de dar unlock

					logger.Warn().
						Str("retry_after", remaining.String()).
						Msg("Comando em cooldown.")

					// Respondemos com um ephemeral para avisar o user
					sendEphemeral(s, i, fmt.Sprintf("Comando em cooldown. Tente novamente em %.1f segundos.", remaining.Seconds()))
					return
				}
			}
			// Finalmente atualizamos o mapa
			userCooldowns[key] = time.Now()
			cooldownMutex.Unlock()
		}

		// Verificamos se o comando quer uma role
		if len(cmd.RequiredRoles) > 0 {
			authorized := false

			// Se no futuro houver atraso na verificação podemos usar um map
			// Double loop para verificar todas as roles
		RoleCheckLoop:
			for _, userRole := range i.Member.Roles {
				for _, requiredRole := range cmd.RequiredRoles {
					if userRole == requiredRole {
						authorized = true
						break RoleCheckLoop
					}
				}
			}

			if !authorized {
				logger.Warn().
					Strs("user_roles", i.Member.Roles).
					Strs("required_roles", cmd.RequiredRoles).
					Msg("Tentativa de comando não autorizado.")

				sendEphemeral(s, i, "Acesso restrito aos cargos autorizados.")
				return
			}
			logger.Debug().Msg("Comando autorizado.")
		}

		// Verificamos se o handler existe
		if cmd.Handler == nil {
			logger.Error().Msg("Comando existe no mapa mas Handler é nulo.")
			sendEphemeral(s, i, "Erro interno do bot.")
			return
		}

		logger.Info().Msg("Comando handler a ser executado.")

		// Executamos o comando
		// cmd.Handler(s, i)

		Enqueue(s, i, cmd.Handler)
	} else {
		// Caso do comando não existir
		logger.Warn().Msg("Interação recebida para comando desconhecido.")
		sendEphemeral(s, i, "Comando desconhecido.")

	}
}

func sendEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Error().Err(err).Str("component", "interactionHandler").Msg("Erro ao responder ao comando (Ephemeral).")
	}
}
