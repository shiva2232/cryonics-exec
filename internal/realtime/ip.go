package realtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SetIP(uid, deviceID, token, ip string) error {
	// Construct the Firebase RTDB path
	url := fmt.Sprintf(
		"https://cryonics-em-default-rtdb.asia-southeast1.firebasedatabase.app/users/%s/%s/ip.json?auth=%s",
		uid, deviceID, token,
	)

	// Prepare the data (must be a JSON string in Firebase)
	data, err := json.Marshal(ip)
	if err != nil {
		return fmt.Errorf("error marshaling IP: %w", err)
	}

	// Send a PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("upload failed, status: %s", resp.Status)
	}

	log.Println("IP address updated successfully")
	return nil
}
