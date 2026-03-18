package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// SendEphemeral envia um ephemeral (generalizado)
func SendEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Erro ao responder (Ephemeral).")
	}
}

// DeferResponse confirma imediatamente a interação para evitar expiração do token
func DeferResponse(s *discordgo.Session, i *discordgo.InteractionCreate, ephemeral bool) error {
	data := &discordgo.InteractionResponseData{}
	if ephemeral {
		data.Flags = discordgo.MessageFlagsEphemeral
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: data,
	})
	if err != nil {
		log.Error().Err(err).Msg("Erro ao deferir resposta.")
	}

	return err
}

// EditDeferredResponse atualiza a resposta já deferida da interação
func EditDeferredResponse(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
	if err != nil {
		log.Error().Err(err).Msg("Erro ao editar resposta deferida.")
	}
}
