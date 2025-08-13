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

func ListenForCommands(ctx context.Context, uid, deviceId string) {
	switch {
	case uid == "":
		log.Println("uid is empty")
	case deviceId == "":
		log.Println("device Id is empty")
	default:
		log.Println("all values are Ok.")
	}
	url := fmt.Sprintf("https://%s.asia-southeast1.firebasedatabase.app/users/%s/%s/commands.json",
		"cryonics-em-default-rtdb", uid, deviceId)

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
				processCommandArray(uid, deviceId, data)
			}
		}
	}
}

func processCommandArray(uid, deviceId, data string) {
	var commands map[string]interface{} // model.CommandData
	if err := json.Unmarshal([]byte(data), &commands); err != nil {
		log.Printf("Failed to parse command array: %v", err)
		log.Println("Failed to data:" + data)
		return
	}
	if commands["path"] == "/" {
		b, err := json.Marshal(commands["data"])
		if err != nil {
			log.Fatal(err)
		}

		var cmds []model.Command
		if err := json.Unmarshal(b, &cmds); err != nil {
			log.Fatal(err)
		}
		for idx, cmd := range cmds {
			if cmd.Status == "pending" {
				log.Println("executing ", cmd.Action)
				go executeAndReport(uid, deviceId, idx, cmd)
			}
		}
	}
}

func executeAndReport(uid, deviceId string, index int, cmd model.Command) {

	updateCommandFields(uid, deviceId, index, map[string]interface{}{
		"status":   "executing",
		"issuedAt": time.Now().Unix(),
	})
	var outputs string = ""
	var errs string = ""

	stdoutCh := make(chan string)
	stderrCh := make(chan string)
	doneCh := make(chan error)

	go executor.ExecuteWithStreaming(cmd.Action, stdoutCh, stderrCh, doneCh)

	for {
		select {
		case out := <-stdoutCh:
			outputs = outputs + `
` + out
			appendCommandOutput(uid, deviceId, index, "[STDOUT] "+outputs)
		case errout := <-stderrCh:
			errs = errs + `
` + errout
			appendCommandOutput(uid, deviceId, index, "[STDERR] "+errs)
		case err := <-doneCh:
			endTime := time.Now()
			status := "completed"
			errorMsg := ""
			if err != nil {
				status = "error"
				errorMsg = err.Error()
			}

			updateCommandFields(uid, deviceId, index, map[string]interface{}{
				"status":      status,
				"completedAt": endTime.Unix(),
				"errorMsg":    errorMsg,
			})
			return
		}
	}
}

func updateCommandFields(uid, deviceId string, index int, fields map[string]interface{}) {
	url := fmt.Sprintf(
		"https://%s.asia-southeast1.firebasedatabase.app/users/%s/%s/commands/%d.json",
		"cryonics-em-default-rtdb", uid, deviceId, index,
	)

	body, _ := json.Marshal(fields)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}

func appendCommandOutput(uid, deviceId string, index int, line string) {
	url := fmt.Sprintf(
		"https://%s.asia-southeast1.firebasedatabase.app/users/%s/%s/commands/%d/output.json",
		"cryonics-em-default-rtdb", uid, deviceId, index,
	)

	body, _ := json.Marshal(line)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}
