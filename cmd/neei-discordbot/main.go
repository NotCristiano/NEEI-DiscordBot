package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"NEEI-DiscordBot/internal/bot"
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"NEEI-DiscordBot/internal/logger"

	"github.com/rs/zerolog/log"
)

func main() {

	// Damos setup no context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Iniciamos o logger
	logger.SetupLogger()
	log.Info().Msg("Iniciando NEEI-DiscordBot...")

	// Carregar config primeiro (antes de init())
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Falha ao carregar config.")
	}
	log.Debug().Msgf("Config carregada.")

	// Registamos todos os comandos encontrados
	commands.InitCommands(cfg)

	// Iniciamos o bot
	discordSession, err := bot.Start(ctx, cfg.Token, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Falha ao iniciar bot.")
	}

	// Esperamos pelo sinal de parada
	<-ctx.Done()

	// Log de intenção shutdown
	log.Info().Msg("A fechar todas as instâncias do bot.")

	// Criamos um timeout para o bot parar caso esteja preso
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Usar uma goroutine para parar o bot, não é necessário, mas garantirmos para futura escalabilidade (DB)
	done := make(chan struct{})
	go func() {
		// Fechar de forma correta os websockets do bot
		if err := discordSession.Close(); err != nil {
			log.Error().Err(err).Msg("Erro ao fechar a sessão do bot.")
		}
		close(done)
	}()

	// Fica à espera aqui até o bot fechar por alguma das duas razões
	select {
	case <-done:
		log.Info().Msg("Sessão do bot fechada.")
	case <-shutdownCtx.Done():
		log.Warn().Msg("Shutdown do bot demorou muito. Parada forçada.")
	}

	// Log do shutdown do bot
	log.Info().Msg("Bot parado.")
}
