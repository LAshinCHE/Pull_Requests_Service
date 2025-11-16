package config

import "os"

type Config struct {
	Addr  string
	DBURL string
}

func LoadConfig() *Config {
	port, _ := os.LookupEnv("PORT")
	dburl, _ := os.LookupEnv("DATABASE_URL")
	return &Config{
		Addr:  port,
		DBURL: dburl,
	}
}
