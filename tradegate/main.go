package main

import (
	"fmt"
	"log"
	"net/http"
	"sidekick/dataservice_client"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var transport *http.Transport

func init() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./log/tradegate.log",
		MaxSize:    100,   // MB
		MaxBackups: 30,    // old files
		MaxAge:     30,    // day
		Compress:   false, // disabled by default
	})
	// log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lmicroseconds | log.Lshortfile)

	transport = &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	}
}

func main() {
	// load config
	config := LoadConfig()

	// new router
	router := mux.NewRouter()

	// init services
	initServices(config, router)

	listenAddr := fmt.Sprintf("%v:%v", config.Host, config.Port)
	log.Printf("listen address: %v\n", listenAddr)
	http.ListenAndServeTLS(listenAddr, "ca/server.crt", "ca/server.key", router)
}

func initServices(config *Config, router *mux.Router) {
	// init data service
	dataService := dataservice_client.NewDataService(config.DataService)

	// init redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB

	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	// init ws handler
	wsHandler := NewWsHandler()
	wsHandler.DataService = dataService
	router.HandleFunc("/v1/ws", wsHandler.Hub.OnAccept)

	// init rest handler
	restHandler := NewRestHandler(router)
	restHandler.DataService = dataService
	restHandler.RedisClient = redisClient
}
