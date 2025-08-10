package utils

import (
	"context"
	"cryonics/internal/iputils"
	"cryonics/internal/metrics"
	"cryonics/internal/realtime"
	"fmt"
	"log"
	"os"
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

func UpdateIpAndMetrics(uid, deviceId, token string) {
	go realtime.ListenForCommands(context.Background(), token, uid, deviceId)

	ip := iputils.GetPublicIP()
	realtime.SetIP(uid, deviceId, token, ip)
	ticker := time.NewTicker(5 * time.Minute)

	defer ticker.Stop()

	for {
		point := metrics.GetMetrics()
		err := realtime.UploadMetric(uid, deviceId, token, point)
		if err != nil {
			log.Println("Upload error:", err)
		}

		<-ticker.C
	}
}
