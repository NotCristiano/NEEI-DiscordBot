package main

import (
	"NEEI-DiscordBot/bot"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {

	// Vamos ler o arquivo .env
	err := godotenv.Load("token.env")

	// Se houver algum erro, o programa vai parar
	if err != nil {
		panic("Erro ao ler o arquivo .env")
	}

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
