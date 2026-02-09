package main

import (
	"NEEI-DiscordBot/bot"
	"NEEI-DiscordBot/commands"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Carregar config primeiro (antes de init())
	commands.LoadConfig()

	// Roles dependem de env e precisam ser definidas depois do LoadConfig.
	commands.SetCommandRoles("echo", []string{commands.RoleDev, commands.RoleDirecao})

	// Extraimos o token do arquivo .env
	token := os.Getenv("TOKEN")

	// Iniciamos o bot
	if err := bot.Start(token); err != nil {
		panic(err)
	}

	// Hook de espera para o bot parar
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Bot parado com sucesso!")
}
