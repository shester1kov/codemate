package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// главная структура конфигурации приложения

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Logger LoggerConfig `mapstructure:"logger"`
	Qdrant QdrantConfig `mapstructure:"qdrant"`
	Ollama OllamaConfig `mapstructure:"ollama"`
}

// настройки сервера

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // режим: debug/release

}

// настройки логгера

type LoggerConfig struct {
	Level      string `mapstructure:"level"`       // уровень логирования: debug/info/warn/error
	Encoding   string `mapstructure:"encoding"`    // формат логов: json/plain
	OutputPath string `mapstructure:"output_path"` // stdout, stderr или путь к файлу
}

// настройки Qdrant
type QdrantConfig struct {
	URL string `mapstructure:"url"`
}

type OllamaConfig struct {
	URL   string `mapstructure:"url"`
	Model string `mapstructure:"model"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("CODEMATE")
	viper.AutomaticEnv()

	setDeafaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil

}

func setDeafaults() {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.encoding", "json")
	viper.SetDefault("logger.output_path", "stdout")

	viper.SetDefault("qdrant.url", "http://localhost:6333")

	viper.SetDefault("ollama.url", "http://localhost:11434")
	viper.SetDefault("ollama.model", "llama3.2:3b")
}
