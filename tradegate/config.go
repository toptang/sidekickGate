package main

import "flag"

type Config struct {
	Host        string
	Port        int
	ServerOkex  string
	DataService string
	Redis       string
}

func LoadConfig() *Config {
	host := flag.String("host", "0.0.0.0", "listen host")
	port := flag.Int("port", 5010, "listen port")
	serverOkex := flag.String("server_okex", "", "okex server")
	dataService := flag.String("data_service", "", "data service")
	redis := flag.String("redis", "", "redis")
	flag.Parse()

	config := &Config{
		Host:        *host,
		Port:        *port,
		ServerOkex:  *serverOkex,
		DataService: *dataService,
		Redis:       *redis,
	}

	return config
}
