package bot

import (
	"NEEI-DiscordBot/commands"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// Start função para iniciar o bot
func Start(token string) (*discordgo.Session, error) {

	// Damos setup no logger e colocamos uma tag do componente
	logger := log.With().Str("component", "bot").Logger()
	logger.Info().Msg("A iniciar o bot...")

	// Iniciamos a queue
	StartQueue()

	// Criamos o bot
	goBot, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error().Err(err).Msg("Falha ao criar sessão do bot.")
		return nil, fmt.Errorf("erro ao criar sessão do bot: %w", err)
	}

	// Chamamos o handler de eventos antes de iniciar o bot
	goBot.AddHandler(interactionHandler)
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
