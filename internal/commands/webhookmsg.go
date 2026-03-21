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
				Name:        "webhookmsg",
				Description: "Envia uma secção de links úteis para um canal específico",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "secao",
						Description: "Secção de links a enviar",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "Links Úteis", Value: "links_uteis"},
							{Name: "Drives de Resumos", Value: "drives-resumos"},
							{Name: "Recursos Genéricos de Informática", Value: "genericos"},
							{Name: "Redes Sociais do NEEI", Value: "redes-sociais"},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "canal",
						Description: "Canal onde a mensagem do webhook será enviada",
						Required:    true,
					},
				},
			},
			Handler:       makeWebhookmsgHandler(cfg),
			RequiredRoles: []string{cfg.RoleDev, cfg.RoleDirecao},
			Cooldown:      5 * time.Second,
			Ephemeral:     true,
		}
	})
}

// makeWebhookmsgHandler retorna a função handler para o comando webhookmsg
func makeWebhookmsgHandler(cfg *config.Config) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		// Trata a logica de opções do comando, visiveis no discord, e guarda os valores em variáveis
		options := i.ApplicationCommandData().Options
		seccao := options[0].StringValue()
		canal := options[1].ChannelValue(s)

		// Verifica se o canal de links está configurado
		if cfg.LinksChannelID == "" {
			EditDeferredResponse(s, i, "❌ Canal de links não configurado. Adiciona `LINKS_CHANNEL_ID` ao `local.env`.")
			return
		}

		// Lê as mensagens do canal de links para extrair as secções e os links (max 20 mensagens)
		messages, err := s.ChannelMessages(cfg.LinksChannelID, 20, "", "", "")
		// Corresponde com feedback de erro caso haja algum problema a ler as mensagens do canal
		if err != nil {
			EditDeferredResponse(s, i, "❌ Erro ao ler o canal de links: "+err.Error())
			return
		}

		// Extrai as secções e os links das mensagens do canal, e procura a secção selecionada pelo utilizador
		sections := ParseLinksFromMessages(messages)
		section, ok := sections[seccao]
		// Corresponde com feedback de erro caso a secção selecionada não seja encontrada
		if !ok {
			EditDeferredResponse(s, i, "❌ Secção `"+seccao+"` não encontrada. Verifica se existe uma mensagem com `SECCAO: "+seccao+"`.")
			return
		}

		// Formata os links da secção para o formato de embed do Discord e cria o embed com a secção e os links e envia para o canal selecionad
		linksValue := ""
		for _, link := range section.Links {
			linksValue += "[" + link.Name + "](" + link.URL + ")\n"
		}

		embed := &discordgo.MessageEmbed{
			Title:       section.Title,
			Description: section.Description,
			Color:       section.Color,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Links", Value: linksValue, Inline: false},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Núcleo de Estudantes de Engenharia Informática da Universidade de Évora",
			},
		}

		// Envia o embed para o canal selecionado e corresponde com feedback de sucesso ou erro
		_, err = s.ChannelMessageSendEmbed(canal.ID, embed)
		if err != nil {
			EditDeferredResponse(s, i, "❌ Erro ao enviar a mensagem: "+err.Error())
			return
		}

		// Feedback de sucesso
		EditDeferredResponse(s, i, "✅ Secção **"+section.Title+"** enviada com sucesso em <#"+canal.ID+">!")
	}
}
