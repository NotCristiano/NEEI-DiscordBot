package main

import (
	"NEEI-DiscordBot/bot"
	"NEEI-DiscordBot/commands"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {

	// Iniciamos o logger
	bot.SetupLogger()
	log.Info().Msg("Iniciando NEEI-DiscordBot...")

	// Carregar config primeiro (antes de init())
	cfg, err := commands.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Falha ao carregar config.")
	}
	log.Debug().Msgf("Config carregada.")

	// Roles dependem de env e precisam ser definidas depois do LoadConfig
	commands.SetCommandRoles("echo", []string{cfg.RoleDev, cfg.RoleDirecao})

	// Iniciamos o bot
	if err := bot.Start(cfg.Token); err != nil {
		log.Fatal().Err(err).Msg("Falha ao iniciar bot.")
	}

	// Hook de espera para o bot parar
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Bloqueia até receber um sinal
	sig := <-sc

	// Log do shutdown do bot
	log.Info().
		Str("signal", sig.String()).
		Msg("Bot parado.")
}
