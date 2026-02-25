package commands

import (
	"NEEI-DiscordBot/internal/config"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func init() {
	AddCommand(func(cfg *config.Config) BotCommand {
		return BotCommand{
			Definition: &discordgo.ApplicationCommand{
				Name:        "timeout",
				Description: "Coloca um membro em timeout",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type: discordgo.ApplicationCommandOptionUser,
						Name: "membro", Description: "Membro a ser colocado em timeout",
						Required: true,
					},
					{
						Type: discordgo.ApplicationCommandOptionString,
						Name: "tempo", Description: "Duração (e.g 1h, 15m, 20s) do timeout",
						Required: true,
					},
				},
			},
			Handler:       timeoutHandler,
			RequiredRoles: []string{cfg.RoleDev, cfg.RoleDirecao},
			Cooldown:      3 * time.Second,
		}
	})
}

// timeoutHandler contém a lógica do comando
func timeoutHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Criamos um logger para o comando
	logger := log.With().Str("component", "commandTimeout").Logger()

	// Extraimos os dados do comando e colocamos num mapa
	options := i.ApplicationCommandData().Options
	optionsMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, option := range options {
		optionsMap[option.Name] = option
	}

	// Recuperamos os dados do comando e verificamos se existe uma razão
	targetUser := optionsMap["membro"].UserValue(s)
	durationStr := optionsMap["tempo"].StringValue()

	// Damos parse no tempo
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		SendEphemeral(s, i, "ERRO: Formato de tempo inválido (e.g 1h, 15m, 20s).")
		logger.Error().Err(err).Str("component", "commandTimeout").Msg("Erro ao formatar tempo.")
		return
	}

	// Verificamos se é maior do que o limite do discord (>24d)
	// https://discordjs.guide/legacy/popular-topics/faq
	if duration.Hours() > 24*28 {
		SendEphemeral(s, i, "ERRO: Não é possível colocar um timeout maior que 28 dias.")
		logger.Warn().Msg("Tentativa de timeout muito grande.")
		return
	}

	// Calculamos até quando o timeout deve ocorrer
	timeoutUntil := time.Now().Add(duration)

	// Editamos o user com GuildMemberEdit
	err = s.GuildMemberTimeout(i.GuildID, targetUser.ID, &timeoutUntil)
	if err != nil {
		SendEphemeral(s, i, "ERRO: Falha ao colocar o membro em timeout.")
		logger.Error().Err(err).Str("component", "commandTimeout").Msg("Erro ao colocar o membro em timeout.")
		return
	}

	// Enviamos uma mensagem ephemeral de sucesso
	SendEphemeral(s, i, "**"+targetUser.Username+"** colocado em timeout com sucesso.")
	logger.Info().Str("component", "commandTimeout").Msg(targetUser.Username + " colocado em timeout por " + durationStr + ".")

}
