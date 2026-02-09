package commands

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// IDs das roles que podem usar os comandos
var (
	RoleDev     string
	RoleDirecao string
)

// LoadConfig carrega o .env antes de qualquer outra coisa
func LoadConfig() {
	err := godotenv.Load("local.env")
	if err != nil {
		fmt.Println("Erro ao carregar arquivo .env:", err)
		return
	}

	RoleDev = os.Getenv("ROLE_DEV")
	RoleDirecao = os.Getenv("ROLE_DIRECAO")
}
