# Guia de Contribuição - NEEI Discord Bot

Este guia existe para manter o trabalho sincronizado e o código limpo, consistente e funcional.

## Quick Start

Se é a tua primeira contribuição, segue estes passos:

1. **Clonar o repositório:**
   ```bash
   git clone https://github.com/NotCristiano/NEEI-DiscordBot.git
   cd NEEI-DiscordBot
   ```

2. **Instalar dependências:**
   ```bash
   go mod download
   ```

3. **Criar o ficheiro `local.env`:**
   O ficheiro deve estar na raiz do projeto e é ignorado pelo Git.

   > Nota: `TOKEN` é obrigatório. As restantes variáveis dependem das funcionalidades que quiseres testar.

   ```env
   TOKEN=TOKEN_DE_TESTE
   SERVER_ID=123456789012345678
   ROLE_DEV=123456789012345678
   ROLE_DIRECAO=123456789012345678
   ROLE_NEEI=123456789012345678
   ROLE_DEPTEC=123456789012345678
   ROLE_DEPAPE=123456789012345678
   ROLE_DEPIMG=123456789012345678
   ROLE_DEPEV=123456789012345678
   ROLE_DEPASSEMBLEIA=123456789012345678
   ROLE_DEPFISCAL=123456789012345678
   FORBIDDEN_CHANNEL_ID=123456789012345678
   AUTO_MUTE_DURATION=10
   ```

4. **Executar o bot:**
   ```bash
   go run cmd/neei-discordbot/main.go
   ```

---

## Como contribuir

A arquitetura do bot é modular. Para criar um comando novo, **normalmente basta criar um ficheiro**.

**Regras:**
1. O ficheiro deve ficar em `internal/commands/`.
2. O package deve ser `package commands`.
3. Tens de usar a função `init()` para registar o comando automaticamente.
4. É obrigatório cada comando ter um teste.

**Exemplo de Comando (Hello):**

```go
package commands

import (
	"NEEI-DiscordBot/internal/config"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Registamos o comando e respetivas restrições
func init() {
	AddCommand(func(cfg *config.Config) BotCommand {
		return BotCommand{
			// Descrição apresentada no Discord
			Definition: &discordgo.ApplicationCommand{
				Name:        "hello",
				Description: "Retorna 'Hello World!'",
				Options:     []*discordgo.ApplicationCommandOption{},
			},
			Handler:       helloHandler,
			RequiredRoles: nil,
			Cooldown:      0 * time.Second,
		}
	})
}

// helloHandler contém a lógica do comando
func helloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Hello World!"}})
}
```

**Exemplo de Teste (`hello_test.go`):**

Testar ligação direta entre bot e Discord é mais complexo, por isso os testes validam a definição e comportamento esperado do comando.

```go
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

	// Verifica se o handler existe
	if commandFound.Handler == nil {
		t.Fatal("O handler do comando é nulo.")
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
```

---

## Comandos de validação

Antes de abrir PR, valida o código localmente:

```bash
go fmt ./...
go test ./...
go vet ./...
```

Se estiveres a trabalhar em Docker, podes também usar os comandos do `Makefile` descritos no `README.md`.

---

## Fluxo de Trabalho

Para manter a `main` estável, seguimos estas regras:

1. **Nunca faças push direto para a `main`.**
2. **Usa o `.gitignore` para evitar exposição de tokens e ficheiros locais.**
3. **Cria uma branch:**
    * Funcionalidade nova: `feat/nome-do-comando` (ex: `feat/ping`)
    * Correção de erro: `fix/nome-do-erro` (ex: `fix/typo-ping`)
4. **Testa localmente:** garante que o bot liga e o comando funciona.
5. **Abre um Pull Request (PR):**
    * No GitHub, abre um PR da tua branch para a `main`.
    * Adiciona uma descrição do que fizeste.
    * Espera pela aprovação do coordenador ou presidente.

---

## Checklist de PR

Antes de submeter, confirma:

- [ ] Código formatado com `go fmt ./...`
- [ ] Testes a passar com `go test ./...`
- [ ] Sem `fmt.Println` (usar `zerolog`)
- [ ] Sem tokens, IDs sensíveis ou segredos no PR
- [ ] Comando novo registado com `init()` e com teste
- [ ] Descrição do PR clara (o que mudou e porquê)

---

## Estilo de Código e Regras

O PR pode ser reprovado se o código não seguir estas regras:

1.  **Formatação:**
	* O código **tem** de estar formatado com `gofmt`.
    * Isto pode ser automatizado na tua IDE de escolha.
    * Ou corre manualmente: `go fmt ./...`

2.  **Logs:**
    * Não uses `fmt.Println`.
    * Usa o logger `zerolog` que já está no projeto.

3.  **Testes:**
	* Antes de enviares, corre `go test ./...` para garantir que não partiste nada.
	* Mesmo com testes automáticos, validar localmente continua a ser boa prática.

