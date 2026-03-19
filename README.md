# NEEI Discord Bot

## Como executar (Localmente)

1. **Instalar dependências:**
   ```bash
   go mod download
   ```

2. **Criar ficheiro `local.env` (exemplo abaixo):**

| Variável | Obrigatória | Descrição |
|----------|-------------|-----------|
| `TOKEN` | Sim | Token do bot no Discord |
| `SERVER_ID` | Não | ID do servidor (guild) |
| `ROLE_DEV` | Não | ID do cargo de desenvolvimento |
| `ROLE_DIRECAO` | Não | ID do cargo de direção |
| `ROLE_NEEI` | Não | ID do cargo NEEI |
| `ROLE_DEPTEC` | Não | ID do cargo do Departamento Tecnológico |
| `ROLE_DEPAPE` | Não | ID do cargo do Departamento de Apoio ao Estudante |
| `ROLE_DEPIMG` | Não | ID do cargo do Departamento de Imagem |
| `ROLE_DEPEV` | Não | ID do cargo do Departamento de Eventos |
| `ROLE_DEPASSEMBLEIA` | Não | ID do cargo de Assembleia |
| `ROLE_DEPFISCAL` | Não | ID do cargo do Conselho Fiscal |
| `FORBIDDEN_CHANNEL_ID` | Não | ID do canal com regras de mute automático |
| `AUTO_MUTE_DURATION` | Não | Duração do mute automático (em minutos) |

Exemplo rápido de `local.env`:

```env
TOKEN=coloca_aqui_o_token
SERVER_ID=123456789012345678
ROLE_DEV=123456789012345678
ROLE_DIRECAO=123456789012345678
ROLE_NEEI=123456789012345678
```


3. **Executar o bot:**
   ```bash
   go run cmd/neei-discordbot/main.go
   ```

---

## Como executar (Docker)

1. **Forma recomendada (Makefile):**
   ```bash
   make up
   ```
   Para parar:
   ```bash
   make down
   ```

2. **Produção (Makefile):**
   ```bash
   make prod-up
   ```
   Para parar:
   ```bash
   make prod-down
   ```
   Para ver logs:
   ```bash
   make prod-logs
   ```

3. **Alternativa sem `make` (manual):**

   **Build da imagem:**
   ```bash
   docker build -t neei-bot .
   ```

4. **Executar o container sem docker-compose:**
   ```bash
   docker run --name neei-bot -v $(pwd)/local.env:/root/local.env neei-bot
   ```

5. **Remover o container caso necessário:**
   ```bash
   docker rm neei-bot
   ```
---

## Estado do Projeto

- [x] Comandos slash base: `hello` e `echo`
- [x] Moderação: `kick`, `ban` e `timeout`
- [x] Controlo de acesso por cargos em comandos sensíveis
- [x] Cooldown por utilizador/comando
- [x] Sistema de tickets (abrir por botão/modal e fechar por botão)
- [x] Auto-mute em canal proibido com duração configurável
- [x] Fila de execução para respeitar rate limits da API do Discord

---

## Contribuir
Para contribuir leia [CONTRIBUTING.md](https://github.com/NotCristiano/NEEI-DiscordBot/blob/master/CONTRIBUTING.md).

---

## Copyright

Copyright (c) the respective contributors, as shown in [Contributors](https://github.com/NotCristiano/NEEI-DiscordBot/graphs/contributors).

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
