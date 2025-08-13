package realtime

import (
	"bytes"
	"cryonics/internal/metrics"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UploadMetric(uid, deviceID string, newMetric metrics.MetricPoint) error {
	basePath := fmt.Sprintf("https://cryonics-em-default-rtdb.asia-southeast1.firebasedatabase.app/users/%s/%s/metrics.json", uid, deviceID)

	// Step 1: Read existing metrics
	resp, err := http.Get(basePath)
	if err != nil {
		return fmt.Errorf("error fetching existing metrics: %w", err)
	}
	defer resp.Body.Close()

	var existing []metrics.MetricPoint
	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 && string(body) != "null" {
			if err := json.Unmarshal(body, &existing); err != nil {
				return fmt.Errorf("error unmarshalling existing metrics: %w", err)
			}
		}
	}

	// Step 2: Append new metric
	existing = append(existing, newMetric)

	// Step 3: Keep last 4
	if len(existing) > 4 {
		existing = existing[len(existing)-4:]
	}

	// Step 4: Upload updated array
	data, err := json.Marshal(existing)
	if err != nil {
		return fmt.Errorf("error marshalling metrics: %w", err)
	}

	req, err := http.NewRequest("PUT", basePath, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	putResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error uploading metrics: %w", err)
	}
	defer putResp.Body.Close()

	if putResp.StatusCode >= 300 {
		return fmt.Errorf("upload failed, status: %s", putResp.Status)
	}
	return nil
}
