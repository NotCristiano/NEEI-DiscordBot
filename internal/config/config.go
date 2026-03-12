package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Token           string `mapstructure:"TOKEN"`
	RoleDev         string `mapstructure:"ROLE_DEV"`
	RoleDirecao     string `mapstructure:"ROLE_DIRECAO"`
	RoleNEEI        string `mapstructure:"ROLE_NEEI"`
	RoleTecnologico string `mapstructure:"ROLE_DEPTEC"`
	RoleApoio       string `mapstructure:"ROLE_DEPAPE"`
	RoleImagem      string `mapstructure:"ROLE_DEPIMG"`
	RoleEventos     string `mapstructure:"ROLE_DEPEV"`
	RoleAssembleia  string `mapstructure:"ROLE_DEPASSEMBLEIA"`
	RoleFiscal      string `mapstructure:"ROLE_DEPFISCAL"`
}

func LoadConfig() (*Config, error) {

	// Especificamos o ficheiro env
	viper.SetConfigName("local")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Caso quisermos usar docker
	viper.SetEnvPrefix("NEEI")
	viper.AutomaticEnv()

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
