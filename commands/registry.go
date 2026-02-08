package commands

import (
	"NEEI-DiscordBot/commands/all"

	"github.com/bwmarrin/discordgo"
)

// HandlerFunc é um Enum para os tipos de comandos
type HandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

// CommandMap mapeia o nome do comando para o seu handler
var CommandMap = map[string]HandlerFunc{
	"hello": all.HelloHandler,
	"echo":  all.EchoHandler,
}

// ComandosApresentar é uma lista para mostrar os comandos disponíveis no discord
var ComandosApresentar = []*discordgo.ApplicationCommand{
	all.HelloCommand,
	all.EchoCommand,
}
