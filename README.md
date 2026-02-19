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
| `ROLE_DEV`    | ID de uma role                           |
| `ROLE_DIRECAO`| ID de uma role                           |
| `ROLE_NEEI`   | ID de uma role                           |


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
