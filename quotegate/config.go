package main

import (
	"flag"
)

type Config struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	CORS_Origin      string `json:"origin"`
	ServerOkex       string `json:"server_okex"`
	ServerSimulation string `json"server_simulation"`
	DataService      string `json:"data_service"`
}

func LoadConfig() *Config {
	host := flag.String("host", "0.0.0.0", "listen host")
	port := flag.Int("port", 5001, "listen port")
	origin := flag.String("CORS_Origin", "", "CORS origin")
	serverOkex := flag.String("server_okex", "", "okex server")
	serverSimulation := flag.String("server_simulation", "", "simulation server")
	dataService := flag.String("data_service", "", "data service")
	flag.Parse()

	config := &Config{
		Host:             *host,
		Port:             *port,
		CORS_Origin:      *origin,
		ServerOkex:       *serverOkex,
		ServerSimulation: *serverSimulation,
		DataService:      *dataService,
	}

	return config
}
