package main

import (
	"fmt"
	"log"
	"net/http"
	"sidekick/dataservice_client"

	"github.com/gorilla/mux"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./log/quotegate.log",
		MaxSize:    100,   // MB
		MaxBackups: 30,    // old files
		MaxAge:     30,    // day
		Compress:   false, // disabled by default
	})
	// log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lmicroseconds | log.Lshortfile)
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
	// init client service
	clientService := NewClientService()
	router.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		clientService.Hub.OnAccept(w, r)
	})

	// init data service
	dataService := dataservice_client.NewDataService(config.DataService)

	// init sub manager
	subManager := NewSubManager()

	// for CORS
	reqHeader := make(map[string][]string)
	if config.CORS_Origin != "" {
		reqHeader["Origin"] = []string{config.CORS_Origin}
	}

	// init okex service
	okexService := NewQuoteService("okex")
	if config.ServerOkex != "" {
		go okexService.Hub.ConnectAndRun(config.ServerOkex, true, 3, reqHeader)
	}

	// Dependency Injection
	clientService.DataService = dataService
	clientService.SubManager = subManager
	clientService.QuoteServices["okex"] = okexService.Hub

	okexService.ClientService = clientService
}
