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
