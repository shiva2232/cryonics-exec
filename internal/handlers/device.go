package handlers

import (
	"cryonics/internal/model"
	"log"
	"net/http"
	"net/url"
)

var DeviceChan = make(chan string)

func ReceiveDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	deviceId := r.FormValue("deviceId")
	name := r.FormValue("name")
	deviceType := r.FormValue("deviceType")
	email := r.FormValue("email") // Must be sent from JS
	uid := r.FormValue("uid")

	if deviceId == "" || name == "" || deviceType == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	dvc := make(map[string]model.Device)
	dvc[deviceId] = model.Device{Name: name, DeviceType: deviceType}

	DeviceChan <- deviceId

	log.Printf("Received device info - ID: %s, Name: %s, Type: %s\n", deviceId, name, deviceType)

	params := url.Values{}
	params.Set("deviceId", deviceId)
	params.Set("deviceName", name)
	params.Set("deviceType", deviceType)
	params.Set("email", email)
	params.Set("uid", uid)

	// Redirect to thank you page
	http.Redirect(w, r, "/thank-you?"+params.Encode(), http.StatusSeeOther)
}
