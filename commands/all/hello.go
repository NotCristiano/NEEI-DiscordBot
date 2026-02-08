package all

import "github.com/bwmarrin/discordgo"

// HelloCommand contém a definição do comando
var HelloCommand = &discordgo.ApplicationCommand{
	Name:        "hello",
	Description: "Retorna 'Hello World!'",
}

// HelloHandler contém a lógica do comando
func HelloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Hello World!"}})
}
