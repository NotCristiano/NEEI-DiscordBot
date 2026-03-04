package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Token       string `mapstructure:"TOKEN"`
	RoleDev     string `mapstructure:"ROLE_DEV"`
	RoleDirecao string `mapstructure:"ROLE_DIRECAO"`
	RoleNEEI    string `mapstructure:"ROLE_NEEI"`
}

func LoadConfig() (*Config, error) {

	// Especificamos o ficheiro env
	viper.SetConfigName("local")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Lê variáveis de ambiente diretamente (TOKEN, ROLE_DEV, etc.)
	viper.AutomaticEnv()

	// Bind explícito para que Unmarshal() consiga ler env vars (necessário para Docker)
	viper.BindEnv("TOKEN")
	viper.BindEnv("ROLE_DEV")
	viper.BindEnv("ROLE_DIRECAO")
	viper.BindEnv("ROLE_NEEI")

	// Ajuda se tivermos nested keys
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Lemos o ficheiro .env
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Warn().Msg("Ficheiro local.env não encontrado.")
		} else {
			return nil, err
		}
	}

	// Desconstrói para uma struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Verificamos campos obrigatórios
	if config.Token == "" {
		log.Fatal().Msg("ERRO: TOKEN vazio; verifique o arquivo local.env")
		return nil, fmt.Errorf("TOKEN vazio no local.env")
	}

	return &config, nil
}
