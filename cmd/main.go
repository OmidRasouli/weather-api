package main

import (
	"log"
	"strconv"

	configs "github.com/OmidRasouli/weather-api/config"
	"github.com/OmidRasouli/weather-api/router"
)

func main() {
	log.Println("The application is starting...")
	cfg := configs.MustLoad("config/config.yaml")
	RunServer(cfg)
}

func RunServer(cfg *configs.Config) {
	// Create the HTTP router and register routes.
	router := router.Setup()

	// Prepare the server address using the configured port.
	port := cfg.GetServerConfig().Port
	addr := ":" + strconv.Itoa(port)

	log.Printf("Server is starting on port %d", port)
	// Start the HTTP server and log any fatal errors.
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
