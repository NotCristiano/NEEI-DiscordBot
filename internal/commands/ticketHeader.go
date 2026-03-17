package commands

import (
	"NEEI-DiscordBot/internal/config"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	AddCommand(func(cfg *config.Config) BotCommand {
		return BotCommand{

			Definition: &discordgo.ApplicationCommand{
				Name:        "ticketheader",
				Description: "Cria a mensagem que contém o botão para criar um ticket.",
				Options:     []*discordgo.ApplicationCommandOption{},
			},
			Handler:       ticketHeaderHandler,
			RequiredRoles: nil,
			Cooldown:      5 * time.Second,
		}

	})
}

func ticketHeaderHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Definimos o embed
	embed := &discordgo.MessageEmbed{
		Title:       "Criar Ticket",
		Description: "Precisas de ajuda, queres sugerir ou reportar alguma coisa? Clica no botão abaixo e entramos em contacto contigo assim que possível! \n\n **Lembra-te que uma resposta pode demorar até 48h!** \n\n ⚠️ Esta funcionalidade é apenas para suporte, por favor não abusem.",
		Color:       0x970302,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Núcleo de Estudantes de Engenharia Informática",
			IconURL: "https://example.com/logo.png",
		},
	}

	// Definimos o botão
	button := discordgo.Button{
		Label:    "Criar Ticket",
		Style:    discordgo.PrimaryButton,
		CustomID: "ticketCreateButton",
		Emoji: &discordgo.ComponentEmoji{
			Name: "💌",
		},
	}

	// Agora juntamos embed e button para enviar
	SendEphemeral(s, i, "Header criado com sucesso!")
	s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{button},
			},
		},
	})

}
