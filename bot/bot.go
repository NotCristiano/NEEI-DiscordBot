package bot

import (
	"NEEI-DiscordBot/commands"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// ServerID os comandos aparecerem mais rápido especificamos o ID do servidor
const ServerID = "1466460166374428786"

// Start função para iniciar o bot
func Start(token string) error {

	// Verificamos redundantemente se o token existe
	if token == "" {
		return fmt.Errorf("TOKEN vazio; verifique o arquivo token.env")
	}

	// Criamos o bot
	goBot, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("erro ao criar sessão do bot: %w", err)
	}

	// Chamamos o handler de eventos antes de iniciar o bot
	goBot.AddHandler(interactionHandler)

	// Se o erro for nulo, o bot foi iniciado com sucesso
	err = goBot.Open()
	if err != nil {
		return fmt.Errorf("erro ao iniciar o bot: %w", err)
	}

	if goBot.State == nil || goBot.State.User == nil {
		return fmt.Errorf("estado do bot não inicializado após Open")
	}

	// Damos load nos comandos
	fmt.Println("carregando comandos...")
	for _, comm := range commands.ComandosApresentar {
		_, err := goBot.ApplicationCommandCreate(goBot.State.User.ID, ServerID, comm)
		if err != nil {
			return fmt.Errorf("erro ao carregar comando: %w", err)
		}
	}

	fmt.Println("bot iniciado com sucesso!")
	return nil
}
