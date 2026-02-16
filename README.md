# NEEI Discord Bot

## Como executar (Localmente)

1. **Instalar dependências:**
   ```bash
   go mod download
   ```

2. **Criar ficheiro `local.env` (exemplo abaixo):**

| Variável      | Descrição                                |
|---------------|------------------------------------------|
| `TOKEN`       | Token do bot Discord                     |
| `ROLE_ID1`    | ID de uma role                           |
| `ROLE_ID2`    | ID de uma role                           |
| `ENVIRONMENT` | Ambiente (`development` ou `production`) |


3. **Executar o bot:**
   ```bash
   go run cmd/neei-discordbot/main.go
   ```

---

## Como executar (Docker)

1. **Build da imagem:**
   ```bash
   docker build -t neei-bot .
   ```
2. **Executar o container sem docker-compose:**
   ```bash
   docker run --name neei-bot -v $(pwd)/local.env:/root/local.env neei-discordbot
   ```
3. **Remover o container caso necessário:**
   ```bash
   docker rm neei-bot
   ```
---

## Contribuir
Para contribuir leia [CONTRIBUTING.md](https://github.com/NotCristiano/NEEI-DiscordBot/blob/master/CONTRIBUTING.md).
