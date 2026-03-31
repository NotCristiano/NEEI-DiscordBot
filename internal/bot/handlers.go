package bot

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"NEEI-DiscordBot/internal/queue"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// maxCooldownDuration é o maior cooldown que qualquer comando pode ter.
// Usado para limpar entradas expiradas do mapa de cooldowns.
const maxCooldownDuration = 10 * time.Minute

// Variáveis globais para dar track no cooldown
var (
	// Formato mapa: UserID: NomeComando -> Ultima execução do comando
	userCooldowns = make(map[string]time.Time)
	cooldownMutex sync.Mutex
)

func interactionHandler(cfg *config.Config) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			handleSlashCommands(s, i)
		case discordgo.InteractionMessageComponent:
			handleButtonInteraction(s, i, cfg)
		case discordgo.InteractionModalSubmit:
			handleModalSubmit(s, i, cfg)
		default:
			log.Warn().Str("interaction_type", string(i.Type)).Msg("Interação recebida de tipo desconhecido. Ignorada.")
		}
	}
}

func handleSlashCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Puxamos as infos do comando
	cmdData := i.ApplicationCommandData()

	// Proteção contra DMs onde i.Member é nulo
	if i.Member == nil || i.Member.User == nil {
		log.Warn().Str("command", cmdData.Name).Msg("Interação recebida sem Member (possível DM). Ignorada.")
		return
	}

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
					commands.SendEphemeral(s, i, fmt.Sprintf("Comando em cooldown. Tente novamente em %.1f segundos.", remaining.Seconds()))
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

				commands.SendEphemeral(s, i, "Acesso restrito aos cargos autorizados.")
				return
			}
			logger.Debug().Msg("Comando autorizado.")
		}

		// Verificamos se o handler existe
		if cmd.Handler == nil {
			logger.Error().Msg("Comando existe no mapa mas Handler é nulo.")
			commands.SendEphemeral(s, i, "Erro interno do bot.")
			return
		}

		logger.Info().Msg("Comando handler a ser executado.")

		// ACK imediato para o token não expirar enquanto o comando é processado
		if err := commands.DeferResponse(s, i, cmd.Ephemeral); err != nil {
			logger.Error().Err(err).Msg("Falha ao defer a resposta da interação.")
			return
		}

		// Damos enqueue no comando
		queue.Enqueue(s, i, cmd.Handler)
	} else {
		// Caso do comando não existir
		logger.Warn().Msg("Interação recebida para comando desconhecido.")
		commands.SendEphemeral(s, i, "Comando desconhecido.")

	}
}

func handleButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, cfg *config.Config) {
	data := i.MessageComponentData()

	switch data.CustomID {
	case "ticketCreateButton":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "ticketCreationModal",
				Title:    "Criar Ticket",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "ticketReason",
								Label:       "Descreva o seu problema:",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Ex. Não consigo falar no canal de voz, ou tenho uma sugestão para o servidor.",
								Required:    true,
								MinLength:   10,
								MaxLength:   2000,
							},
						},
					},
				},
			},
		})
	case "ticketCloseButton":
		_, err := s.ChannelDelete(i.ChannelID)
		if err != nil {
			log.Error().Err(err).Msg("Erro ao fechar ticket.")
			commands.SendEphemeral(s, i, "Erro interno do bot.")
			return
		}

		log.Info().Str("channel_id", i.ChannelID).Msg("Ticket fechado.")
	default:
		// Verifica se é um botão de navegação do comando links
		if strings.HasPrefix(data.CustomID, "links_prev_") || strings.HasPrefix(data.CustomID, "links_next_") {
			handleLinksNavigation(s, i, cfg)
			return
		}
		log.Warn().Str("custom_id", data.CustomID).Msg("Interação de botão recebida com CustomID desconhecido. Ignorada.")
		commands.SendEphemeral(s, i, "Ação desconhecida.")
	}
}

func handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate, cfg *config.Config) {
	data := i.ModalSubmitData()

	switch data.CustomID {
	case "ticketCreationModal":
		if err := commands.DeferResponse(s, i, true); err != nil {
			return
		}

		description := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		channel, err := s.Channel(i.ChannelID)
		if err != nil {
			log.Error().Err(err).Msg("Erro ao obter canal para ticket.")
			commands.EditDeferredResponse(s, i, "Erro interno do bot.")
			return
		}

		createTicket(s, i, description, channel.ParentID, cfg)

	default:
		log.Warn().Str("custom_id", data.CustomID).Msg("Interação de modal recebida com CustomID desconhecido. Ignorada.")
		commands.SendEphemeral(s, i, "Ação desconhecida.")
	}
}

// startCooldownCleaner inicia uma goroutine que periodicamente limpa entradas expiradas do mapa de cooldowns
func startCooldownCleaner(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				cooldownMutex.Lock()
				for key, lastExec := range userCooldowns {
					if now.Sub(lastExec) > maxCooldownDuration {
						delete(userCooldowns, key)
					}
				}
				cooldownMutex.Unlock()
			}
		}
	}()
}

func handleLinksNavigation(s *discordgo.Session, i *discordgo.InteractionCreate, cfg *config.Config) {
	logger := log.With().Str("component", "linksNavigation").Str("user", i.Member.User.Username).Logger()

	data := i.MessageComponentData()

	// Extrai o índice atual do CustomID (ex: "links_next_0" -> 0)
	var currentIndex int
	if strings.HasPrefix(data.CustomID, "links_next_") {
		fmt.Sscanf(data.CustomID, "links_next_%d", &currentIndex)
		currentIndex++
	} else {
		fmt.Sscanf(data.CustomID, "links_prev_%d", &currentIndex)
		currentIndex--
	}

	// Lê as mensagens do canal de links
	messages, err := s.ChannelMessages(cfg.LinksChannelID, 20, "", "", "")
	if err != nil {
		logger.Error().Err(err).Msg("Erro ao ler o canal de links na navegação.")
		commands.SendEphemeral(s, i, "❌ Erro ao ler o canal de links.")
		return
	}

	sections := commands.ParseLinksFromMessages(messages)
	orderedKeys := commands.GetOrderedSectionKeys(sections)

	// Garante que o índice está dentro dos limites
	if currentIndex < 0 {
		currentIndex = 0
	}
	if currentIndex >= len(orderedKeys) {
		currentIndex = len(orderedKeys) - 1
	}

	embed := commands.BuildLinksEmbed(sections, orderedKeys, currentIndex)
	components := commands.BuildLinksComponents(orderedKeys, currentIndex)

	// Atualiza a mensagem existente
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("Erro ao atualizar embed de links.")
		return
	}

	logger.Info().Int("index", currentIndex).Msg("Navegação de links executada com sucesso.")
}
