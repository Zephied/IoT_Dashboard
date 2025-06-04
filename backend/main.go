package main

import (
	"fmt"
	"log"
	"net/http"

	"iot-dashboard/backend/api"
	"iot-dashboard/backend/db"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("[INFO] Initializing database")
	db.Init()

	// Démarre le scan périodique
	go api.StartPeriodicScan()

	fmt.Println("[INFO] Starting IoT Dashboard Server on :8080")
	r := mux.NewRouter()

	r.HandleFunc("/api/devices/{id}/action", api.ControlDeviceHandler).Methods("POST")
	r.HandleFunc("/api/devices/{id}", api.DeleteDeviceHandler).Methods("DELETE")
	r.HandleFunc("/api/devices/{id}", api.UpdateDeviceHandler).Methods("PUT")
	r.HandleFunc("/api/scan", api.ScanHandler)
	r.HandleFunc("/api/devices", api.GetDevicesHandler)
	r.HandleFunc("/api/add-mock-device", api.AddMockDeviceHandler)

	// Sert le frontend (fichiers statiques)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend")))

	log.Fatal(http.ListenAndServe(":8080", r))
}
