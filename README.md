# NEEI Discord Bot

## Como executar

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
   go run main.go
   ```