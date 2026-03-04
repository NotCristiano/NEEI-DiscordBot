package queue

import (
	"NEEI-DiscordBot/internal/commands"
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Request representa o request da execução de um comando
type Request struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Tamanho do buffer da fila de requests
var requestQueue = make(chan Request, 100)

// StartQueue inicia a execução de request que estão na fila
// Aceita um context para permitir shutdown
func StartQueue(ctx context.Context) {
	// Para não exceder a api do discord limitamos a 50 requests por segundo
	// https://docs.discord.com/developers/topics/rate-limits#global-rate-limit:~:text=Global%20Rate%20Limit,-All
	ticker := time.NewTicker(50 * time.Millisecond)

	go func() {
		defer ticker.Stop()
		log.Info().Msg("Queue de requests iniciada.")

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Queue de requests a encerrar.")
				return
			case <-ticker.C:
				// Esperamos por um request na fila ou pelo contexto cancelar
				select {
				case <-ctx.Done():
					log.Info().Msg("Queue de requests a encerrar.")
					return
				case request := <-requestQueue:
					// Executamos os requests
					func() {
						defer func() {
							if r := recover(); r != nil {
								log.Error().Str("panic", fmt.Sprintf("%v", r)).Msg("Erro ao executar request na fila.")
							}
						}()
						request.Handler(request.Session, request.Interaction)
					}()
				}
			}
		}
	}()
}

// Enqueue adiciona um request na fila
func Enqueue(s *discordgo.Session, i *discordgo.InteractionCreate, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	select {
	case requestQueue <- Request{s, i, handler}:
		// Foi adicionado no buffer
		log.Debug().Msg("Request adicionado na fila.")
	default:
		// Buffer cheio
		log.Warn().Msg("Buffer de requests cheio.")
		commands.SendEphemeral(s, i, "Bot sobrecarregado. Tente novamente em alguns segundos.")
	}
}
