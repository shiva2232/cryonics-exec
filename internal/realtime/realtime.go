package realtime

import (
	"cryonics/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchUserDevices(idToken, uid string) (model.Devices, error) {
	url := fmt.Sprintf("https://cryonics-em-default-rtdb.asia-southeast1.firebasedatabase.app/users/%s.json?auth=%s", uid, idToken)

	resp, err := http.Get(url)
	if err != nil {
		return model.Devices{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return model.Devices{}, fmt.Errorf("error from RTDB: %s - %s", resp.Status, string(body))
	}

	var data model.Devices
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return model.Devices{}, err
	}
	return data, nil
}
