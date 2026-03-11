package tests

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TestEchoCommand(t *testing.T) {
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
		if cmd.Definition.Name == "echo" {
			commandFound = cmd
			found = true
			break
		}
	}

	// Verifica se o comando foi encontrado
	if !found {
		t.Fatal("O comando 'echo' não foi encontrado.")
	}

	// Verifica a Descrição
	expectedDesc := "Repete o que é inserido"
	if commandFound.Definition.Description != expectedDesc {
		t.Fatalf("Descrição incorreta do comando 'echo'. Esperado: %s, Encontrado: %s ", expectedDesc, commandFound.Definition.Description)
	}

	// Verifica o numero de opções
	if len(commandFound.Definition.Options) != 1 {
		t.Fatalf("Esperava uma opção no comando echo, mas encontrou %d", len(commandFound.Definition.Options))
	}

	// Verifica as opções
	opt := commandFound.Definition.Options[0]
	if opt.Name != "texto" {
		t.Fatalf("Nome da opção incorreto. Esperado: 'texto', Encontrado: '%s'", opt.Name)
	}

	// Verifica o tipo de opções
	if opt.Type != discordgo.ApplicationCommandOptionString {
		t.Fatalf("Tipo de opção incorreto. Esperado: String, Encontrado %d", opt.Type)
	}

	// Verifica a obrigatoriedade
	if !opt.Required {
		t.Fatalf("A opção 'texto' deve ser obrigatória. Esperado: true, Encontrado: %t", opt.Required)
	}

	// Verifica se o Handler existe
	if commandFound.Handler == nil {
		t.Fatal("O Handler do comando é nulo.")
	}

	// Verifica o cooldown
	if commandFound.Cooldown != 5*time.Second {
		t.Fatalf("O cooldown devia ser 5, mas é: %s ", commandFound.Cooldown)
	}

	// Verifica as roles requeridas
	expectedRole := []string{"roledev", "roledirecao"}

	// Verifica o numero de roles
	if len(commandFound.RequiredRoles) != len(expectedRole) {
		t.Fatalf("Número de roles incorretos. Esperado: %d, Encontrado: %d", len(expectedRole), len(commandFound.RequiredRoles))
	}

	// Verifica se as roles corretas foram adicionadas
	for i, role := range commandFound.RequiredRoles {
		if role != expectedRole[i] {
			t.Fatalf("Role incorreta. Esperado: %s, Encontrado: %s", expectedRole[i], role)
		}
	}
}
