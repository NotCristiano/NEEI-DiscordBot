package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Token                   string `mapstructure:"TOKEN"`
	ServerID                string `mapstructure:"SERVER_ID"`
	RoleDev                 string `mapstructure:"ROLE_DEV"`
	RoleDirecao             string `mapstructure:"ROLE_DIRECAO"`
	RoleNEEI                string `mapstructure:"ROLE_NEEI"`
	RoleTecnologico         string `mapstructure:"ROLE_DEPTEC"`
	RoleApoio               string `mapstructure:"ROLE_DEPAPE"`
	RoleImagem              string `mapstructure:"ROLE_DEPIMG"`
	RoleEventos             string `mapstructure:"ROLE_DEPEV"`
	RoleAssembleia          string `mapstructure:"ROLE_DEPASSEMBLEIA"`
	RoleFiscal              string `mapstructure:"ROLE_DEPFISCAL"`
	RolePresiDeptec         string `mapstructure:"ROLE_PRESI_DEPTEC"`
	RolePresiApe            string `mapstructure:"ROLE_PRESI_APE"`
	RolePresiImg            string `mapstructure:"ROLE_PRESI_IMG"`
	RolePresiEv             string `mapstructure:"ROLE_PRESI_EV"`
	RolePresiAssembleia     string `mapstructure:"ROLE_PRESI_ASSEMBLEIA"`
	RolePresiFiscal         string `mapstructure:"ROLE_PRESI_FISCAL"`
	RoleVicePresiDeptec     string `mapstructure:"ROLE_VICE_PRESI_DEPTEC"`
	RoleVicePresiApe        string `mapstructure:"ROLE_VICE_PRESI_APE"`
	RoleVicePresiImg        string `mapstructure:"ROLE_VICE_PRESI_IMG"`
	RoleVicePresiEv         string `mapstructure:"ROLE_VICE_PRESI_EV"`
	RoleVicePresiAssembleia string `mapstructure:"ROLE_VICE_PRESI_ASSEMBLEIA"`
	RoleVicePresiFiscal     string `mapstructure:"ROLE_VICE_PRESI_FISCAL"`
	MuteChannelID           string `mapstructure:"FORBIDDEN_CHANNEL_ID"`
	AutoMuteDuration        int    `mapstructure:"AUTO_MUTE_DURATION"`
	LinksChannelID          string `mapstructure:"LINKS_CHANNEL_ID"`
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
	viper.BindEnv("SERVER_ID")
	viper.BindEnv("ROLE_DEV")
	viper.BindEnv("ROLE_DIRECAO")
	viper.BindEnv("ROLE_NEEI")
	viper.BindEnv("ROLE_DEPTEC")
	viper.BindEnv("ROLE_DEPAPE")
	viper.BindEnv("ROLE_DEPIMG")
	viper.BindEnv("ROLE_DEPEV")
	viper.BindEnv("ROLE_DEPASSEMBLEIA")
	viper.BindEnv("ROLE_DEPFISCAL")
	viper.BindEnv("ROLE_PRESI_DEPTEC")
	viper.BindEnv("ROLE_PRESI_APE")
	viper.BindEnv("ROLE_PRESI_IMG")
	viper.BindEnv("ROLE_PRESI_EV")
	viper.BindEnv("ROLE_PRESI_ASSEMBLEIA")
	viper.BindEnv("ROLE_PRESI_FISCAL")
	viper.BindEnv("ROLE_VICE_PRESI_DEPTEC")
	viper.BindEnv("ROLE_VICE_PRESI_APE")
	viper.BindEnv("ROLE_VICE_PRESI_IMG")
	viper.BindEnv("ROLE_VICE_PRESI_EV")
	viper.BindEnv("ROLE_VICE_PRESI_ASSEMBLEIA")
	viper.BindEnv("ROLE_VICE_PRESI_FISCAL")
	viper.BindEnv("FORBIDDEN_CHANNEL_ID")
	viper.BindEnv("AUTO_MUTE_DURATION")
	viper.BindEnv("LINKS_CHANNEL_ID")

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
