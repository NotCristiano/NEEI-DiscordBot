package tests

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"testing"
	"time"
)

func TestLinksCommand(t *testing.T) {
	// Criamos uma config dummy
	mockConfig := &config.Config{
		RoleDev:     "roledev",
		RoleDirecao: "roledirecao",
	}

	// Buscamos o comando
	var commandFound commands.BotCommand
	found := false

	// Iteramos sobre os comandos
	for _, cmdConfig := range commands.ComandosLista {
		cmd := cmdConfig(mockConfig)
		if cmd.Definition.Name == "links" {
			commandFound = cmd
			found = true
			break
		}
	}

	// Verifica se o comando foi encontrado
	if !found {
		t.Fatal("O comando 'links' não foi encontrado.")
	}

	// Verifica a descrição
	if commandFound.Definition.Description != "Mostra as secções de links úteis do NEEI-UÉ" {
		t.Fatalf("Descrição incorreta: %s", commandFound.Definition.Description)
	}

	// Verifica se o Handler existe
	if commandFound.Handler == nil {
		t.Fatal("O handler do comando é nulo.")
	}

	// Verifica que não tem roles requeridas (qualquer membro pode usar)
	if len(commandFound.RequiredRoles) != 0 {
		t.Fatalf("O comando links não deve ter roles requeridas, mas tem: %v", commandFound.RequiredRoles)
	}

	// Verifica que não tem opções (comando sem argumentos)
	if len(commandFound.Definition.Options) != 0 {
		t.Fatalf("O comando links não deve ter opções, mas tem: %d", len(commandFound.Definition.Options))
	}

	// Verifica o cooldown
	if commandFound.Cooldown != 5*time.Second {
		t.Fatalf("O cooldown devia ser 5s, mas é: %s", commandFound.Cooldown)
	}

	// Verifica que é ephemeral
	if !commandFound.Ephemeral {
		t.Fatal("O comando links devia ser ephemeral.")
	}
}

func TestGetOrderedSectionKeys(t *testing.T) {
	// Criamos secções de teste
	sections := map[string]commands.LinkSection{
		"links-uteis":    {Title: "Links Úteis"},
		"drives-resumos": {Title: "Drives de Resumos"},
		"genericos":      {Title: "Recursos Genéricos"},
	}

	keys := commands.GetOrderedSectionKeys(sections)

	// Verifica que a ordem está correta
	expectedOrder := []string{"links-uteis", "drives-resumos", "genericos"}
	if len(keys) != len(expectedOrder) {
		t.Fatalf("Esperado %d keys, mas tem: %d", len(expectedOrder), len(keys))
	}

	// Verifica que as keys estão na ordem correta
	for idx, key := range expectedOrder {
		if keys[idx] != key {
			t.Fatalf("Ordem incorreta na posição %d. Esperado: %s, Encontrado: %s", idx, key, keys[idx])
		}
	}
}
