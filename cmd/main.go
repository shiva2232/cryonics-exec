package main

import (
	"cryonics/internal/utils"
	"log"
)

func main() {
	uid, device, err := utils.DecryptFile()
	if err != nil {
		log.Println("invalid encrypted file or file not found")
	}

	utils.UpdateIpAndMetrics(uid, device)
}
