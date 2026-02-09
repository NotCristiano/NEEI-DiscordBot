package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
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
	CommandMap[cmd.Definition.Name] = cmd
	ComandosApresentar = append(ComandosApresentar, cmd.Definition)
}

// SetCommandRoles atualiza as roles requeridas depois de carregar config.
func SetCommandRoles(name string, roles []string) {
	cmd, ok := CommandMap[name]
	if !ok {
		fmt.Println("Comando não encontrado para atualizar roles:", name)
		return
	}
	cmd.RequiredRoles = roles
	CommandMap[name] = cmd
}
