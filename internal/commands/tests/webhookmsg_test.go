package tests

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"testing"
	"time"
)

func TestWebhookmsgCommand(t *testing.T) {
	// Criamos uma config dummy
	mockConfig := &config.Config{
		RoleDev:         "roledev",
		RoleDirecao:     "roledirecao",
		RolePresiDeptec: "rolepresideptec",
	}

	// Buscamos o comando
	var commandFound commands.BotCommand
	found := false

	// Iteramos sobre os comandos
	for _, cmdConfig := range commands.ComandosLista {
		cmd := cmdConfig(mockConfig)
		if cmd.Definition.Name == "webhookmsg" {
			commandFound = cmd
			found = true
			break
		}
	}

	// Verifica se o comando foi encontrado
	if !found {
		t.Fatal("O comando 'webhookmsg' não foi encontrado.")
	}

	// Verifica a Descrição
	if commandFound.Definition.Description != "Envia uma secção de links úteis para um canal específico" {
		t.Fatalf("Descrição incorreta: %s", commandFound.Definition.Description)
	}

	// Verifica se o Handler existe
	if commandFound.Handler == nil {
		t.Fatal("O handler do comando é nulo.")
	}

	// Verifica as roles requeridas
	expectedRole := []string{"roledev", "roledirecao", "rolepresideptec"}
	if len(commandFound.RequiredRoles) != len(expectedRole) {
		t.Fatalf("Devia ter %d roles requeridas, mas tem: %d", len(expectedRole), len(commandFound.RequiredRoles))
	}
	for i, role := range commandFound.RequiredRoles {
		if role != expectedRole[i] {
			t.Fatalf("Role incorreta. Esperado: %s, Encontrado: %s", expectedRole[i], role)
		}
	}

	// Verifica o número de opções
	if len(commandFound.Definition.Options) != 2 {
		t.Fatalf("Devia ter 2 opções, mas tem: %d", len(commandFound.Definition.Options))
	}

	// Verifica o nome da primeira opção
	if commandFound.Definition.Options[0].Name != "secao" {
		t.Fatalf("Primeira opção devia ser 'secao', mas é: %s", commandFound.Definition.Options[0].Name)
	}

	// Verifica o nome da segunda opção
	if commandFound.Definition.Options[1].Name != "canal" {
		t.Fatalf("Segunda opção devia ser 'canal', mas é: %s", commandFound.Definition.Options[1].Name)
	}

	// Verifica o número de choices da primeira opção
	if len(commandFound.Definition.Options[0].Choices) != 4 {
		t.Fatalf("Devia ter 4 choices, mas tem: %d", len(commandFound.Definition.Options[0].Choices))
	}

	// Verifica o cooldown
	if commandFound.Cooldown != 5*time.Second {
		t.Fatalf("O cooldown devia ser 5s, mas é: %s", commandFound.Cooldown)
	}
}

// TestWebhookLinkSections verifica se as secções de links estão corretamente definidas e têm os campos necessários
func TestWebhookLinkSections(t *testing.T) {
	expectedSections := []string{"links-uteis", "drives-resumos", "genericos", "redes-sociais"}

	// Verifica se todas as secções esperadas estão presentes e têm título e descrição
	for _, key := range expectedSections {
		section, ok := commands.SectionMeta[key]
		if !ok {
			t.Fatalf("Secção '%s' não encontrada no SectionMeta.", key)
		}
		if section.Title == "" {
			t.Fatalf("Secção '%s' não tem título.", key)
		}
		if section.Description == "" {
			t.Fatalf("Secção '%s' não tem descrição.", key)
		}
		// Links não são verificados aqui pois vêm dinamicamente do canal Discord
	}
}
