package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// BotCommand é a estrutura para os comandos dos bots
type BotCommand struct {
	Definition    *discordgo.ApplicationCommand // Extrai tudo do comando
	Handler       func(s *discordgo.Session, i *discordgo.InteractionCreate)
	RequiredRoles []string
	Cooldown      time.Duration
}

// CommandConfig é um type que precisa da config para retornar um comando
type CommandConfig func(cfg *Config) BotCommand

var (
	// CommandMap mapeia o nome do comando para o seu handler
	CommandMap = make(map[string]BotCommand)

	// ComandosApresentar é uma lista para mostrar os comandos disponíveis no discord
	ComandosApresentar []*discordgo.ApplicationCommand

	// ComandosLista é uma lista que contém todos os comandos para registar após o load da config
	ComandosLista []CommandConfig
)

// AddCommand adiciona um novo comando à lista de comandos para registar
func AddCommand(cmd CommandConfig) {
	ComandosLista = append(ComandosLista, cmd)
}

func InitCommands(cfg *Config) {

	logger := log.With().Str("component", "initCommands").Logger()

	logger.Debug().Int("count", len(ComandosLista)).Msg("Registando comandos.")

	// Iteramos sobre a lista de comandos
	for _, cmd := range ComandosLista {
		// Registamos o comando
		RegisterCommand(cmd(cfg))
	}
}

// RegisterCommand regista um novo comando no bot
func RegisterCommand(cmd BotCommand) {

	// Criamos um logger para o comando
	logger := log.With().Str("component", "registerCommand").Logger()

	// Verificamos se os dados são nulos
	if cmd.Definition == nil {
		logger.Error().Msg("Comando com dados nulos tentou ser registado.")
		return
	}

	// Agora podemos ter um logger com o nome do comando
	cmdLogger := logger.With().Str("command", cmd.Definition.Name).Logger()

	// Verificamos duplicidade de comando, edge case
	if _, exists := CommandMap[cmd.Definition.Name]; exists {
		cmdLogger.Warn().Msg("Tentativa de registar comando duplicado.")
	}

	// Registamos o comando
	CommandMap[cmd.Definition.Name] = cmd
	ComandosApresentar = append(ComandosApresentar, cmd.Definition)

	// Log de sucesso
	cmdLogger.Debug().
		Str("description", cmd.Definition.Description).
		Int("options_count", len(cmd.Definition.Options)).
		Msg("Comando registado com sucesso.")
}
