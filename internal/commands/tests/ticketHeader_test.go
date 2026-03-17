package tests

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"testing"
	"time"
)

func TestTicketHeaderCommand(t *testing.T) {
	// Criamos uma config dummy
	mockConfig := &config.Config{
		RoleDev:     "roledev",
		RoleDirecao: "roledirecao",
		RoleNEEI:    "roleneei",
	}

	// Buscamos o comando
	var commandFound commands.BotCommand
	found := false

	// Iteramos sobre os comandos
	for _, cmdConfig := range commands.ComandosLista {
		cmd := cmdConfig(mockConfig)
		if cmd.Definition.Name == "ticketheader" {
			commandFound = cmd
			found = true
			break
		}
	}

	// Verifica se o comando foi encontrado
	if !found {
		t.Fatal("O comando 'ticketheader' não foi encontrado.")
	}

	// Verifica a Descrição
	expectedDesc := "Cria a mensagem que contém o botão para criar um ticket."
	if commandFound.Definition.Description != expectedDesc {
		t.Fatalf("Descrição incorreta do comando 'ticketheader'. Esperado: %s, Encontrado: %s ", expectedDesc, commandFound.Definition.Description)
	}

	// Verifica o numero de opções
	if len(commandFound.Definition.Options) != 0 {
		t.Fatalf("Esperava zero opções no comando ticketheader, mas encontrou %d", len(commandFound.Definition.Options))
	}

	// Verifica se o Handler existe
	if commandFound.Handler == nil {
		t.Fatal("O Handler do comando é nulo.")
	}

	// Verifica o cooldown
	if commandFound.Cooldown != 5*time.Second {
		t.Fatalf("O cooldown devia ser 5 segundos, mas é: %s ", commandFound.Cooldown)
	}

	// Verifica as roles requeridas (ticketheader não tem roles requeridas)
	if len(commandFound.RequiredRoles) != 0 {
		t.Fatalf("O comando ticketheader não deve ter roles requeridas, mas tem: %v", commandFound.RequiredRoles)
	}

}
