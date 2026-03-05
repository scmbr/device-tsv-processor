package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	EnvLocal = "local"
)

type Config struct {
	Environment string `mapstructure:"environment"`
	Postgres    PostgresConfig
	Rabbit      RabbitConfig
	HTTP        HTTPConfig

	BatchSize   int    `mapstructure:"batch_size"`
	OutputDir   string `mapstructure:"output_dir"`
	MaxAttempts int    `mapstructure:"max_attempts"`
}

type PostgresConfig struct {
	Username string
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string
	SSLMode  string `mapstructure:"sslmode"`
	Password string
}

type HTTPConfig struct {
	Host               string        `mapstructure:"host"`
	Port               int           `mapstructure:"port"`
	ReadTimeout        time.Duration `mapstructure:"readTimeout"`
	WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
	MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
}
type RabbitConfig struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	Username          string
	Password          string
	VHost             string
	FileQueueName     string `mapstructure:"file_queue_name"`
	DocumentQueueName string `mapstructure:"document_queue_name"`
}

func Init(configsDir string) (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}
	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.SetDefaults()
	cfg.OverrideFromEnv()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) SetDefaults() {
	if c.Postgres.SSLMode == "" {
		c.Postgres.SSLMode = "disable"
	}
	if c.HTTP.Port == 0 {
		c.HTTP.Port = 8080
	}
	if c.HTTP.ReadTimeout <= 0 {
		c.HTTP.ReadTimeout = 5 * time.Second
	}
	if c.HTTP.WriteTimeout <= 0 {
		c.HTTP.WriteTimeout = 5 * time.Second
	}
	if c.HTTP.MaxHeaderMegabytes <= 0 {
		c.HTTP.MaxHeaderMegabytes = 1 << 20
	}

}

func (c *Config) OverrideFromEnv() {
	if val := os.Getenv("POSTGRES_USER"); val != "" {
		c.Postgres.Username = val
	}
	if val := os.Getenv("POSTGRES_PASSWORD"); val != "" {
		c.Postgres.Password = val
	}
	if val := os.Getenv("POSTGRES_DB"); val != "" {
		c.Postgres.Name = val
	}
	if val := os.Getenv("RABBIT_USERNAME"); val != "" {
		c.Rabbit.Username = val
	}
	if val := os.Getenv("RABBIT_PASSWORD"); val != "" {
		c.Rabbit.Password = val
	}
	if val := os.Getenv("RABBIT_VHOST"); val != "" {
		c.Rabbit.VHost = val
	}

}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read main config: %w", err)
	}

	if env == "" || env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)
	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("failed to read %s config: %w", env, err)
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Environment == "" {
		return fmt.Errorf("environment is not set")
	}

	if c.Postgres.Host == "" {
		return fmt.Errorf("postgres host is required")
	}
	if c.Postgres.Port <= 0 || c.Postgres.Port > 65535 {
		return fmt.Errorf("postgres port must be between 1 and 65535")
	}
	if c.Postgres.Username == "" {
		return fmt.Errorf("postgres username is required")
	}
	if c.Postgres.Password == "" {
		return fmt.Errorf("postgres password is required")
	}
	if c.Postgres.Name == "" {
		return fmt.Errorf("postgres database name is required")
	}
	if c.Postgres.SSLMode == "" {
		return fmt.Errorf("postgres sslmode is required")
	}
	if c.Rabbit.Host == "" {
		return fmt.Errorf("rabbit host is required")
	}
	if c.Rabbit.Port <= 0 || c.Rabbit.Port > 65535 {
		return fmt.Errorf("rabbit port must be between 1 and 65535")
	}
	if c.Rabbit.Username == "" {
		return fmt.Errorf("rabbit username is required")
	}
	if c.Rabbit.Password == "" {
		return fmt.Errorf("rabbit password is required")
	}
	if c.Rabbit.VHost == "" {
		return fmt.Errorf("rabbit vhost is required")
	}
	if c.Rabbit.FileQueueName == "" {
		return fmt.Errorf("rabbit file_queue_name is required")
	}
	if c.Rabbit.DocumentQueueName == "" {
		return fmt.Errorf("rabbit document_queue_name is required")
	}

	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		return fmt.Errorf("http port must be between 1 and 65535")
	}
	if c.HTTP.ReadTimeout <= 0 || c.HTTP.WriteTimeout <= 0 {
		return fmt.Errorf("http timeouts must be > 0")
	}
	if c.HTTP.MaxHeaderMegabytes <= 0 {
		return fmt.Errorf("http maxHeaderMegabytes must be > 0")
	}
	if c.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}
	if c.BatchSize > 500 || c.BatchSize <= 0 {
		return fmt.Errorf("batch size must be between 0 and 500")
	}
	if c.BatchSize > 10 || c.BatchSize <= 0 {
		return fmt.Errorf("batch size must be between 0 and 10")
	}
	return nil
}
