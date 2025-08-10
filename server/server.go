package server

import (
	"cryonics/internal/handlers"
	"log"
	"net/http"
)

func RunServer(port string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.ServeIndex)
	mux.HandleFunc("/select-device", handlers.ServeSelection)
	mux.HandleFunc("/receive-device", handlers.ReceiveDeviceHandler)
	mux.HandleFunc("/receive-token", handlers.ReceiveToken)
	mux.HandleFunc("/thank-you", handlers.ServeThankYou)
	// Configure server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	go func() {
		log.Printf("Server running at http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	return server
}
