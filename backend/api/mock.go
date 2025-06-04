package api

import (
	"encoding/json"
	"net/http"

	"iot-dashboard/backend/db"
	"iot-dashboard/backend/models"
)

func AddMockDeviceHandler(w http.ResponseWriter, r *http.Request) {
	mock := models.Device{
		IP:   "127.0.0.2",         // ip fictive pour la caméra locale
		MAC:  "00:11:22:33:44:55", // adresse MAC fictive
		Name: "Camera USB Locale",
		Type: "camera",
		Actions: `{
			"voir la caméra": {
				"url": "http://{{ip}}/video",
				"method": "GET"
			}
		}`,
		Online: true,
		IsMock: true,
	}
	db.InsertDevice(mock)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mock)
}
