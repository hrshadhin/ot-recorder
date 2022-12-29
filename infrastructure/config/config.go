package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Hook     HooksConfig    `mapstructure:"hook"`
}

// AppConfig app specific config
type AppConfig struct {
	Env              string        `mapstructure:"env"`
	DataPath         string        `mapstructure:"data_path"`
	RequestBodyLimit string        `mapstructure:"request_body_limit"`
	Host             string        `mapstructure:"host"`
	TimeZone         string        `mapstructure:"time_zone"`
	ReadTimeout      time.Duration `mapstructure:"read_timeout"`
	WriteTimeout     time.Duration `mapstructure:"write_timeout"`
	IdleTimeout      time.Duration `mapstructure:"idle_timeout"`
	ContextTimeout   time.Duration `mapstructure:"context_timeout"`
	Port             int           `mapstructure:"port"`
	Debug            bool          `mapstructure:"debug"`
}

// DatabaseConfig DB specific config
type DatabaseConfig struct {
	Type        string        `mapstructure:"type"`
	Host        string        `mapstructure:"host"`
	Name        string        `mapstructure:"name"`
	Username    string        `mapstructure:"username"`
	Password    string        `mapstructure:"password"`
	SslMode     string        `mapstructure:"ssl_mode"`
	MaxLifeTime time.Duration `mapstructure:"max_life_time"`
	Port        int           `mapstructure:"port"`
	MaxOpenConn int           `mapstructure:"max_open_conn"`
	MaxIdleConn int           `mapstructure:"max_idle_conn"`
	Debug       bool          `mapstructure:"debug"`
}

type HooksConfig struct {
	Telegram TelegramHook `mapstructure:"telegram"`
}

type TelegramHook struct {
	SecretToken string `mapstructure:"secret_token"`
	ChatID      int64  `mapstructure:"chat_id"`
}

// c is the configuration instance
var c Config //nolint:gochecknoglobals

// Get returns all configurations
func Get() Config {
	return c
}

// LoadTestValues Load test config
func LoadTestValues() {
	c.Hook.Telegram.SecretToken = "secret"
	c.Hook.Telegram.ChatID = 1
}

// Load the config
func Load(path string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return fmt.Errorf("failed to unmarshal consul config: %w", err)
	}

	if c.App.RequestBodyLimit == "" {
		c.App.RequestBodyLimit = "1M"
	}

	dataPath := strings.TrimSpace(c.App.DataPath)
	if dataPath == "" {
		dataPath = "."
	}

	c.App.DataPath = dataPath

	if c.App.TimeZone == "" {
		c.App.TimeZone = "UTC"
	}

	return nil
}
