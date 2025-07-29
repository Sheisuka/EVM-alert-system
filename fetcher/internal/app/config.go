package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file was found")
	}
}

type Config struct {
	wssAddress string
}

func NewConfig() *Config {
	return &Config{
		wssAddress: getEnvAsString("WSS_ADDRESS", ""),
	}
}

// WSSAddress returns the Websocket endpoint URL, or a default one if none is provided
func (c *Config) WSSAddress() string {
	return c.wssAddress
}

func getEnvAsString(key, defultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defultValue
}
