package bot

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"NEEI-DiscordBot/internal/queue"
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Start função para iniciar o bot
func Start(ctx context.Context, token string, cfg *config.Config) (*discordgo.Session, error) {

	// Damos setup no logger e colocamos uma tag do componente
	logger := log.With().Str("component", "bot").Logger()
	logger.Info().Msg("A iniciar o bot...")

	// Iniciamos a queue com suporte a shutdown
	queue.StartQueue(ctx)

	// Iniciamos o cleaner de cooldowns
	startCooldownCleaner(ctx)

	// Criamos o bot
	goBot, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error().Err(err).Msg("Falha ao criar sessão do bot.")
		return nil, fmt.Errorf("erro ao criar sessão do bot: %w", err)
	}

	// Chamamos o handler de eventos antes de iniciar o bot
	goBot.AddHandler(interactionHandler(cfg))
	goBot.AddHandler(messageCreateHandler(cfg))
	logger.Debug().Msg("Handler de eventos registado.")

	// Se o erro for nulo, o bot foi iniciado com sucesso
	err = goBot.Open()
	if err != nil {
		logger.Error().Err(err).Msg("Falha ao iniciar websocket do bot.")
		return nil, fmt.Errorf("erro ao iniciar o bot: %w", err)
	}

	// Verificações de estado para garantir que o bot foi iniciado com sucesso
	if goBot.State == nil || goBot.State.User == nil {
		logger.Error().Msg("Estado do bot é nulo mesmo depois de Open.")
		return nil, fmt.Errorf("estado do bot não inicializado após Open")
	}

	// Guardamos para o logger as informações do bot assim que iniciado
	logger.Info().
		Str("nome", goBot.State.User.Username).
		Str("id", goBot.State.User.ID).
		Msg("Bot iniciado com sucesso!")

	// Limpamos comandos antigos que já não existem no registo local
	cleanStaleCommands(goBot, logger)

	// Damos load nos comandos
	logger.Info().Int("count", len(commands.ComandosApresentar)).Msg("Carregando comandos...")
	for _, comm := range commands.ComandosApresentar {

		// Usamos outro subLogger para cada comando
		commLogger := logger.With().Str("command", comm.Name).Logger()

		_, err := goBot.ApplicationCommandCreate(goBot.State.User.ID, "", comm)
		if err != nil {
			commLogger.Error().Err(err).Msg("Falha ao carregar comando.")
			return nil, fmt.Errorf("erro ao carregar comando: %w", err)
		}

		commLogger.Info().Msg("Comando carregado com sucesso!")
	}

	return goBot, nil
}

// cleanStaleCommands remove do Discord comandos que já não existem no registo local
func cleanStaleCommands(s *discordgo.Session, logger zerolog.Logger) {
	registeredCmds, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		logger.Warn().Err(err).Msg("Não foi possível listar comandos do Discord para limpeza.")
		return
	}

	for _, cmd := range registeredCmds {
		if _, exists := commands.CommandMap[cmd.Name]; !exists {
			logger.Info().Str("command", cmd.Name).Msg("A remover comando obsoleto do Discord.")
			if err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID); err != nil {
				logger.Warn().Err(err).Str("command", cmd.Name).Msg("Falha ao remover comando obsoleto.")
			}
		}
	}
}
