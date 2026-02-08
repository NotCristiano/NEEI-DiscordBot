package commands

import (
	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "hello",
		Description: "Retorna 'Hello World!'",
	},
}
