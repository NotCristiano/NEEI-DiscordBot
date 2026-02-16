# Guia de Contribuição - NEEI Discord Bot

Este guia serve para garantir que todos trabalhem de forma sincronizada e que o código se mantenha limpo e funcional para os próximos mandátos.

## Como Começar

Se é a primeira vez que vais contribuir, seguem os passos:

1.  **Clonar o Repositório:**
    ```bash
    git clone https://github.com/NotCristiano/NEEI-DiscordBot.git
    cd NEEI-DiscordBot
    ```

2.  **Configurar Variáveis Locais:**
    Cria um ficheiro chamado `local.env` na raiz do projeto (este ficheiro é ignorado pelo Git).

    Dados usádos não serão revelados, mas contém a seguinte estrutura:
    ```env
    TOKEN=TOKEN_DE_TESTE
    ROLE_DEV=123456789
    ROLE_DIRECAO=987654321
    ROLE_NEEI=1122334455
    ```

3.  **Executar o Bot:**
    ```bash
    go run cmd/neei-discordbot/main.go
    ```

---

## Como contribuir

A arquitetura do bot é modular. Para criar um comando novo, **só é necessário criar um ficheiro**.

**Regras:**
1. O ficheiro deve ir para: `internal/commands/`.
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

// Automaticamente registamos o comando e especificamos os dados e restrições
func init() {
	AddCommand(func(cfg *config.Config) BotCommand {
		return BotCommand{

			// Descrição que vai aparecer no discord
			Definition: &discordgo.ApplicationCommand{
				Name:        "hello",
				Description: "Retorna 'Hello World!'",
				Options:     []*discordgo.ApplicationCommandOption{},
			},
			Handler:       HelloHandler, // Função que vai ser executada quando o comando for executado
			RequiredRoles: nil, // Roles que precisam ter para executar o comando (i.e []string{cfg.RoleDev, cfg.RoleDirecao})
			Cooldown:      0 * time.Second, // Cooldown para o comando
		}
	})
}

// HelloHandler contém a lógica do comando
func HelloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Hello World!"}})
}
```
---

## Fluxo de Trabalho

Para manter a `main` estável, seguimos estas regras:

1.  **Nunca será feito um push direto na `main`.**
2. **Usar o `.gitignore` fornecido para não dar leak em tokens.**
3. **Cria uma Branch:**
    * Funcionalidade nova: `feat/nome-do-comando` (ex: `feat/ping`)
    * Correção de erro: `fix/nome-do-erro` (ex: `fix/typo-ping`)
4. **Testa Localmente:** Garante que o bot liga e o comando funciona.
5. **Abre um Pull Request (PR):**
    * No GitHub, abre um PR da tua branch para a `main`.
    * Adiciona uma descrição do que fizeste.
    * Espera pela aprovação do coordenador ou presidente.

---

## Estilo de Código e Regras

O PR será reprovado se o código se não seguir isto:

1.  **Formatação:**
    * O código **tem** de estar formatado com `gofmt`.
    * Isto pode ser automatizado na tua IDE de escolha.
    * Ou corre manualmente: `go fmt ./...`

2.  **Logs:**
    * Não uses `fmt.Println`.
    * Usa o logger `zerolog` que já está no projeto.

3.  **Testes:**
    * Antes de enviares, corre `go test ./...` para garantir que não partiste nada. Mesmo que existam testes automáticos é boa prática.

