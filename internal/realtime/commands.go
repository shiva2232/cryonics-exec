package realtime

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"cryonics/internal/executor"
	"cryonics/internal/model"
)

func ListenForCommands(ctx context.Context, token, uid, deviceId string) {
	switch {
	case token == "":
		log.Println("token is empty")
	case uid == "":
		log.Println("uid is empty")
	case deviceId == "":
		log.Println("device Id is empty")
	default:
		log.Println("all values are Ok.")
	}
	url := fmt.Sprintf("https://%s.asia-southeast1.firebasedatabase.app/users/%s/%s/commands.json?auth=%s",
		"cryonics-em-default-rtdb", uid, deviceId, token)
	log.Println(url, "is url")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to connect to Firebase: %v", err)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	log.Println(resp.Body)
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping command listener...")
			return
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Println(err, "is error")
				if err.Error() != "EOF" {
					log.Printf("Read error: %v", err)
				}
				time.Sleep(2 * time.Second)
				continue
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Firebase sends "event:" and "data:"
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimPrefix(line, "data:")
				data = strings.TrimSpace(data)
				if data == "null" || data == "" {
					continue
				}
				processCommandArray(token, uid, deviceId, data)
			}
		}
	}
}

func processCommandArray(token, uid, deviceId, data string) {
	var commands model.CommandData
	if err := json.Unmarshal([]byte(data), &commands); err != nil {
		log.Printf("Failed to parse command array: %v", err)
		log.Println("Failed to parse command array:", data)
		return
	}
	log.Println(commands, commands.Data, "in mid")
	for idx, cmd := range commands.Data {
		if cmd.Status == "pending" {
			log.Printf("Pending command found: -> %s", cmd.Action)
			go executeAndReport(token, uid, deviceId, idx, cmd)
		}
	}
}

func executeAndReport(token, uid, deviceId string, index int, cmd model.Command) {

	updateCommandFields(token, uid, deviceId, index, map[string]interface{}{
		"status":   "executing",
		"issuedAt": time.Now().Unix(),
	})

	stdoutCh := make(chan string)
	stderrCh := make(chan string)
	doneCh := make(chan error)

	go executor.ExecuteWithStreaming(cmd.Action, stdoutCh, stderrCh, doneCh)

	for {
		select {
		case out := <-stdoutCh:
			appendCommandOutput(token, uid, deviceId, index, "[STDOUT] "+out)
		case errout := <-stderrCh:
			appendCommandOutput(token, uid, deviceId, index, "[STDERR] "+errout)
		case err := <-doneCh:
			endTime := time.Now()
			status := "completed"
			errorMsg := ""
			if err != nil {
				status = "error"
				errorMsg = err.Error()
			}

			updateCommandFields(token, uid, deviceId, index, map[string]interface{}{
				"status":      status,
				"completedAt": endTime.Unix(),
				"errorMsg":    errorMsg,
			})
			return
		}
	}
}

func updateCommandFields(token, uid, deviceId string, index int, fields map[string]interface{}) {
	url := fmt.Sprintf(
		"https://%s.asia-southeast1.firebasedatabase.app/users/%s/%s/commands/%d.json?auth=%s",
		"cryonics-em-default-rtdb", uid, deviceId, index, token,
	)

	body, _ := json.Marshal(fields)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}

func appendCommandOutput(token, uid, deviceId string, index int, line string) {
	url := fmt.Sprintf(
		"https://%s.asia-southeast1.firebasedatabase.app/users/%s/%s/commands/%d/output.json?auth=%s",
		"cryonics-em-default-rtdb", uid, deviceId, index, token,
	)

	body, _ := json.Marshal(line)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}
