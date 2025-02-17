package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"time"
)

// Config Структура конфигурации;
// Содержит все конфигурационные данные о сервисе;
// автоподгружается при изменении исходного файла
type Config struct {
	ServiceHost string
	ServicePort int

	Redis RedisConfig
}

// RedisConfig Структура конфигурации Redis
type RedisConfig struct {
	Host        string
	Port        int
	Password    string
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

// NewConfig Создаёт новый объект конфигурации, загружая данные из файла конфигурации
func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	// анмаршал конфига
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	cfg.Redis.Port, err = strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Error("error parsing redis port: %w", err)
	}
	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Redis.User = os.Getenv("REDIS_USER")
	log.Info("config parsed")

	return cfg, nil
}
