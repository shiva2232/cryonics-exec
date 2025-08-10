package main

import (
	"context"
	"cryonics/internal/auth"
	"cryonics/internal/handlers"
	"cryonics/internal/utils"
	"cryonics/server"
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		// log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	uid := os.Getenv("UID")
	device := os.Getenv("DEVICEID")

	var reauth = flag.Bool("nocache", false, "forces to reauthentication")
	flag.Parse()

	if _, err := auth.VerifyFirebaseIDToken(context.Background(), token); token == "" || uid == "" || device == "" || err != nil || *reauth {
		port := "8080"
		log.Println("visit localhost:" + port + "/ to run setup")
		server := server.RunServer(port)

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
		utils.UpdateIpAndMetrics(usr.UID, deviceId, usr.Token)
	} else {
		utils.UpdateIpAndMetrics(uid, device, token)
	}
}
