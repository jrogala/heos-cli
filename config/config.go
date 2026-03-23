package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	DefaultPort = 1255
)

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func Init() {
	viper.SetDefault("port", DefaultPort)

	viper.SetEnvPrefix("HEOS")
	viper.BindEnv("host")
	viper.BindEnv("port")

	configDir, err := os.UserConfigDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(configDir, "heos-cli"))
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	_ = viper.ReadInConfig()
}

func Get() Config {
	return Config{
		Host: viper.GetString("host"),
		Port: viper.GetInt("port"),
	}
}

func ConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(configDir, "heos-cli")
}

func Save(host string, port int) error {
	dir := ConfigDir()
	if dir == "" {
		return fmt.Errorf("cannot determine config directory")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	viper.Set("host", host)
	viper.Set("port", port)

	return viper.WriteConfigAs(filepath.Join(dir, "config.yaml"))
}
