package tests

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"testing"
	"time"
)

func TestHelloCommand(t *testing.T) {
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
		if cmd.Definition.Name == "hello" {
			commandFound = cmd
			found = true
			break
		}
	}

	// Verifica se o comando foi encontrado
	if !found {
		t.Fatal("O comando 'hello' não foi encontrado.")
	}

	// Verifica a Descrição
	expectedDesc := "Retorna 'Hello World!'"
	if commandFound.Definition.Description != expectedDesc {
		t.Fatalf("Descrição incorreta do comando 'hello'. Esperado: %s, Encontrado: %s ", expectedDesc, commandFound.Definition.Description)
	}

	// Verifica se o Handler existe
	if commandFound.Handler == nil {
		t.Fatal("O Handler do comando é nulo.")
	}

	// Verifica as roles requeridas (neste caso hello não tem)
	if len(commandFound.RequiredRoles) != 0 {
		t.Fatalf("O comando Hello não deve ter roles requeridas, mas tem: %v", commandFound.RequiredRoles)
	}

	// Verifica o cooldown
	if commandFound.Cooldown != 0*time.Second {
		t.Fatalf("O cooldown devia ser 0, mas é: %s ", commandFound.Cooldown)
	}

}
