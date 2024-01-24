package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

var Mode, Port, WebPrefix string
var AllowedOrigins []string

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func init() {

	// Mode
	var argMode string
	flag.StringVar(&argMode, "APP_MODE", "", "")
	flag.Parse()

	if len(argMode) > 0 {
		// command line argument overrides environment variable
		Mode = argMode
		log.Printf("[main] APP_MODE from command line: %v\n", Mode)
	} else if envMode := os.Getenv("APP_MODE"); len(envMode) > 0 {
		Mode = envMode
		log.Printf("[main] APP_MODE from environment variable: %v\n", Mode)
	}

	// Port
	Port = getenv("APP_PORT", "8100")

	// WebPrefix
	WebPrefix = getenv("APP_WEB_PREFIX", "")

	// Origins
	envOrigins := os.Getenv("APP_ORIGINS")
	if len(envOrigins) > 0 {
		AllowedOrigins = append(AllowedOrigins, strings.Split(envOrigins, ",")...)
	}
	if Mode == "DEV" || len(envOrigins) == 0 {
		AllowedOrigins = append(AllowedOrigins, "http://localhost:8100", "https://localhost:8100")
	}
}
