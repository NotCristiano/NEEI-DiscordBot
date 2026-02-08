package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Start função para iniciar o bot
func Start(token string) error {

	if token == "" {
		return fmt.Errorf("TOKEN vazio; verifique o arquivo token.env")
	}

	goBot, err := discordgo.New("Bot " + token) // Criamos o bot
	if err != nil {
		return fmt.Errorf("erro ao criar sessão do bot: %w", err)
	}

	// Se o erro for nulo, o bot foi iniciado com sucesso
	err = goBot.Open()

	if err != nil {
		return fmt.Errorf("erro ao iniciar o bot: %w", err)
	}

	if goBot.State == nil || goBot.State.User == nil {
		return fmt.Errorf("estado do bot não inicializado após Open")
	}

	fmt.Println("bot iniciado com sucesso!")
	return nil
}
