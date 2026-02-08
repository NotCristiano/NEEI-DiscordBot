package all

import "github.com/bwmarrin/discordgo"

// EchoCommand contém a descrição do comando e o argumento necessário
var EchoCommand = &discordgo.ApplicationCommand{
	Name:        "echo",
	Description: "Repete o que é inserido",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "texto", Description: "Texto a ser repetido",
			Required: true,
		},
	},
}

// EchoHandler contém a lógica do comando
func EchoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Extraimos o texto do comando
	options := i.ApplicationCommandData().Options

	// A mensagem a ser repetida vai estar na primeira posição
	msg := options[0].StringValue()

	// Output do comando
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	})
}
