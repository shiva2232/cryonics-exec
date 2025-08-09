package main

import (
	"context"
	"cryonics/internal/handlers"
	"cryonics/internal/realtime"
	"log"
	"net/http"
	"time"
)

func main() {
	keep := make(chan string, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.ServeIndex)
	mux.HandleFunc("/select-device", handlers.ServeSelection)
	mux.HandleFunc("/receive-device", handlers.ReceiveDeviceHandler)
	mux.HandleFunc("/receive-token", handlers.ReceiveToken)
	mux.HandleFunc("/thank-you", handlers.ServeThankYou)

	// Configure server
	port := "8080"
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

	usr := <-handlers.UserChan
	_, err := realtime.FetchUserDevices(usr.Token, usr.UID)
	if err != nil {
		log.Println("Error:", err)
	}
	deviceId := <-handlers.DeviceChan
	<-handlers.End
	<-time.After(1 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	realtime.ListenForCommands(context.Background(), usr.Token, usr.UID, deviceId)

	<-keep
}
