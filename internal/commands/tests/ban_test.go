package tests

import (
	"NEEI-DiscordBot/internal/commands"
	"NEEI-DiscordBot/internal/config"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestCommandBan(t *testing.T) {
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
		if cmd.Definition.Name == "ban" {
			commandFound = cmd
			found = true
			break
		}
	}

	// Verifica se o comando foi encontrado
	if !found {
		t.Fatal("O comando 'ban' não foi encontrado.")
	}

	// Verifica a Descrição
	expectedDesc := "Banimento de um user especifico"
	if commandFound.Definition.Description != expectedDesc {
		t.Fatalf("Descrição incorreta do comando 'ban'. Esperado: %s, Encontrado: %s ", expectedDesc, commandFound.Definition.Description)
	}

	// Verifica o numero de opções (2)
	if len(commandFound.Definition.Options) != 3 {
		t.Fatalf("Esperava três opções no comando 'ban', mas encontrou %d", len(commandFound.Definition.Options))
	}

	// Verifica a opção
	opt1 := commandFound.Definition.Options[0]
	if opt1.Name != "membro" {
		t.Fatalf("Nome da opção incorreto. Esperado: 'membro', Encontrado: '%s'", opt1.Name)
	}

	// Verifica o tipo de opções
	if opt1.Type != discordgo.ApplicationCommandOptionUser {
		t.Fatalf("Tipo de opção incorreto. Esperado: User, Encontrado %d", opt1.Type)
	}

	// Verifica a obrigatoriedade
	if !opt1.Required {
		t.Fatalf("A opção 'membro' deve ser obrigatória. Esperado: true, Encontrado: %t", opt1.Required)
	}

	opt2 := commandFound.Definition.Options[1]
	if opt2.Name != "dias" {
		t.Fatalf("Nome da opção incorreto. Esperado: 'dias', Encontrado: '%s'", opt2.Name)
	}

	// Verifica o tipo de opções
	if opt2.Type != discordgo.ApplicationCommandOptionInteger {
		t.Fatalf("Tipo de opção incorreto. Esperado: Number, Encontrado %d", opt2.Type)
	}

	// Verifica a obrigatoriedade
	if !opt2.Required {
		t.Fatalf("A opção 'dias' deve ser obrigatória. Esperado: true, Encontrado: %t", opt2.Required)
	}

	opt3 := commandFound.Definition.Options[2]
	if opt3.Name != "razao" {
		t.Fatalf("Nome da opção incorreto. Esperado: 'razão', Encontrado: '%s'", opt3.Name)
	}

	// Verifica o tipo de opções
	if opt3.Type != discordgo.ApplicationCommandOptionString {
		t.Fatalf("Tipo de opção incorreto. Esperado: String, Encontrado %d", opt3.Type)
	}

	// Verifica a obrigatoriedade
	if !opt3.Required {
		t.Fatalf("A opção 'razão' deve ser obrigatória. Esperado: true, Encontrado: %t", opt3.Required)
	}

	// Verifica se o Handler existe
	if commandFound.Handler == nil {
		t.Fatal("O Handler do comando é nulo.")
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
