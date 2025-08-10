package utils

import (
	"fmt"
	"log"
	"os"
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
