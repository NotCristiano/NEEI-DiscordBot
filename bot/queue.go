package bot

import (
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
func StartQueue() {
	// Para não exceder a api do discord limitamos a 50 requests por segundo
	// https://docs.discord.com/developers/topics/rate-limits#global-rate-limit:~:text=Global%20Rate%20Limit,-All
	ticker := time.NewTicker(50 * time.Millisecond)

	go func() {
		log.Info().Msg("Queue de requests iniciada.")
		for {
			// Esperamos pelo ticker
			<-ticker.C

			// Esperamos por um request na fila
			request := <-requestQueue

			// Executamos os requests
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Error().Err(err.(error)).Msg("Erro ao executar request na fila.")
					}
				}()
				request.Handler(request.Session, request.Interaction)
			}()
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
		sendEphemeral(s, i, "Bot sobrecarregado. Tente novamente em alguns segundos.")
	}
}
