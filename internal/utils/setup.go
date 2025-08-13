package utils

import (
	"context"
	"cryonics/internal/iputils"
	"cryonics/internal/metrics"
	"cryonics/internal/realtime"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func CreateEnv(usrToken, usrUID, deviceId string) {
	// Create or overwrite .env file
	file, err := os.Create(".env")
	if err != nil {
		log.Printf("Error creating .env file: %v", err)
	}
	defer file.Close()

	// Write environment variables
	_, err = file.WriteString(fmt.Sprintf("TOKEN=%s\nUID=%s\nDEVICEID=%s\n", usrToken, usrUID, deviceId))
	if err != nil {
		log.Printf("Error writing to .env file: %v", err)
	}

	log.Println(".env file created successfully")
}

func UpdateIpAndMetrics(uid, deviceId string) {
	go realtime.ListenForCommands(context.Background(), uid, deviceId)

	ip := iputils.GetPublicIP()
	realtime.SetIP(uid, deviceId, ip)
	if runtime.GOOS != "android" {
		ticker := time.NewTicker(5 * time.Minute)

		defer ticker.Stop()

		for {
			point := metrics.GetMetrics()
			err := realtime.UploadMetric(uid, deviceId, point)
			if err != nil {
				log.Println("Upload error:", err)
			}

			<-ticker.C
		}
	}
}
