package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// BotCommand é a estrutura para os comandos dos bots
type BotCommand struct {
	Definition    *discordgo.ApplicationCommand // Extrai tudo do comando
	Handler       func(s *discordgo.Session, i *discordgo.InteractionCreate)
	RequiredRoles []string
}

var (
	// CommandMap mapeia o nome do comando para o seu handler
	CommandMap = make(map[string]BotCommand)

	// ComandosApresentar é uma lista para mostrar os comandos disponíveis no discord
	ComandosApresentar []*discordgo.ApplicationCommand
)

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

// SetCommandRoles atualiza as roles requeridas depois de carregar config
func SetCommandRoles(name string, roles []string) {

	// Setup do logger
	logger := log.With().
		Str("component", "setCommandRoles").
		Str("command", name).
		Logger()

	cmd, ok := CommandMap[name]

	// Verificamos se o comando existe
	if !ok {
		logger.Error().
			Strs("attempted_roles", roles).
			Msg("Tentativa de atualizar roles de comando inexistente.")
		return
	}
	cmd.RequiredRoles = roles
	CommandMap[name] = cmd

	// Log de sucesso
	logger.Debug().
		Int("roles_count", len(roles)).
		Strs("roles", roles).
		Msg("Roles atualizados com sucesso.")
}
