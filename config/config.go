package config

import (
	"log"

	"github.com/spf13/viper"
)



type Config interface {	
	Get(key string) interface{} 
	GetString(key string) string
	IsSet(key string) bool
	GetInt(key string) int64
}

type config struct {
	cfg *viper.Viper
}

func New() Config {
	
	cfg := *viper.New()
	cfg.SetConfigName("config.yaml")
	cfg.SetConfigType("yaml")
	cfg.AddConfigPath("./")
	cfg.AddConfigPath("./config")

	if err := cfg.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	return &config{cfg: &cfg}
}


func (c *config) Get(key string) interface{} {
	return c.cfg.Get(key)
}

func (c *config) GetString(key string) string {
	return c.cfg.GetString(key)
}

func (c *config) IsSet(key string) bool {
	return c.cfg.IsSet(key)
}

func (c *config) GetInt(key string) int64 {
	return c.cfg.GetInt64(key)
}
