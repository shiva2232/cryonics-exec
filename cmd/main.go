package main

import (
	"context"
	"cryonics/internal/auth"
	"cryonics/internal/handlers"
	"cryonics/internal/realtime"
	"cryonics/internal/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	keep := make(chan string, 1)
	err := godotenv.Load(".env")
	if err != nil {
		// log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	uid := os.Getenv("UID")
	device := os.Getenv("DEVICEID")

	if _, err := auth.VerifyFirebaseIDToken(context.Background(), token); token == "" || uid == "" || device == "" || err != nil {
		port := "8080"
		log.Println("visit localhost:" + port + "/ to run setup")
		server := runServer(port)

		usr := <-handlers.UserChan
		// _, err = realtime.FetchUserDevices(usr.Token, usr.UID)
		// if err != nil {
		// 	log.Println("Error:", err)
		// }
		deviceId := <-handlers.DeviceChan
		<-handlers.End
		<-time.After(1 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		utils.CreateEnv(usr.Token, usr.UID, deviceId)
		realtime.ListenForCommands(context.Background(), usr.Token, usr.UID, deviceId)

	} else {
		realtime.ListenForCommands(context.Background(), token, uid, device)
	}

	<-keep
}

func runServer(port string) *http.Server {
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
