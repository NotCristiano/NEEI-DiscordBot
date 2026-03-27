package commands

import (
	"NEEI-DiscordBot/internal/config"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func init() {
	AddCommand(func(cfg *config.Config) BotCommand {
		return BotCommand{
			Definition: &discordgo.ApplicationCommand{
				Name:        "links",
				Description: "Mostra as secções de links úteis do NEEI-UÉ",
			},
			Handler:       makeLinksHandler(cfg),
			RequiredRoles: nil,
			Cooldown:      5 * time.Second,
			Ephemeral:     true,
		}
	})
}

// makeLinksHandler retorna o handler para o comando links
func makeLinksHandler(cfg *config.Config) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		logger := log.With().Str("command", "links").Str("user", i.Member.User.Username).Logger()

		// Verifica se o canal de links está configurado
		if cfg.LinksChannelID == "" {
			logger.Warn().Msg("LINKS_CHANNEL_ID não configurado.")
			EditDeferredResponse(s, i, "❌ Canal de links não configurado.")
			return
		}

		// Lê as mensagens do canal de links para extrair as secções e os links
		messages, err := s.ChannelMessages(cfg.LinksChannelID, 20, "", "", "")
		if err != nil {
			logger.Error().Err(err).Msg("Erro ao ler o canal de links.")
			EditDeferredResponse(s, i, "❌ Erro ao ler o canal de links: "+err.Error())
			return
		}

		// Extrai as secções dos links a partir das mensagens
		sections := ParseLinksFromMessages(messages)
		if len(sections) == 0 {
			logger.Warn().Msg("Nenhuma secção encontrada no canal de links.")
			EditDeferredResponse(s, i, "❌ Nenhuma secção encontrada. Verifica o canal de links.")
			return
		}

		// Ordena as secções pela ordem definida no SectionMeta
		orderedKeys := GetOrderedSectionKeys(sections)

		// Envia a primeira secção com botões de navegação
		embed := BuildLinksEmbed(sections, orderedKeys, 0)
		components := BuildLinksComponents(orderedKeys, 0)

		// Edita a resposta da interação para mostrar o embed e os botões
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds:     &[]*discordgo.MessageEmbed{embed},
			Components: &components,
		})
		if err != nil {
			logger.Error().Err(err).Msg("Erro ao enviar embed de links.")
			return
		}

		logger.Info().Msg("Comando links executado com sucesso.")
	}
}

// GetOrderedSectionKeys devolve as chaves das secções na ordem definida no SectionMeta
func GetOrderedSectionKeys(sections map[string]LinkSection) []string {
	// Ordem fixa das secções
	order := []string{"links-uteis", "drives-resumos", "genericos", "redes-sociais"}
	var keys []string
	for _, key := range order {
		if _, ok := sections[key]; ok {
			keys = append(keys, key)
		}
	}
	return keys
}

// BuildLinksEmbed constrói o embed para a secção no índice dado
func BuildLinksEmbed(sections map[string]LinkSection, keys []string, index int) *discordgo.MessageEmbed {
	section := sections[keys[index]]

	// Formata os links da secção para o formato de embed do Discord
	linksValue := ""
	for _, link := range section.Links {
		linksValue += "[" + link.Name + "](" + link.URL + ")\n"
	}

	// Cria o embed com a secção e os links
	return &discordgo.MessageEmbed{
		Title:       section.Title,
		Description: section.Description,
		Color:       section.Color,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Links", Value: linksValue, Inline: false},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Secção %d de %d — Núcleo de Estudantes de Engenharia Informática da Universidade de Évora", index+1, len(keys)),
		},
	}
}

// BuildLinksComponents constrói os botões de navegação
func BuildLinksComponents(keys []string, index int) []discordgo.MessageComponent {
	// Determina se os botões de navegação devem estar desativados (ex: não pode ir para a secção anterior se estiver na primeira ou para a seguinte se estiver na última)
	prevDisabled := index == 0
	nextDisabled := index == len(keys)-1

	// Cria os botões de navegação com custom IDs que indicam a ação e o índice da secção
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "⬅️ Anterior",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("links_prev_%d", index),
					Disabled: prevDisabled,
				},
				discordgo.Button{
					Label:    "➡️ Seguinte",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("links_next_%d", index),
					Disabled: nextDisabled,
				},
			},
		},
	}
}
