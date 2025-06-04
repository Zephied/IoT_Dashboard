package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"iot-dashboard/backend/db"
	"iot-dashboard/backend/scanner"

	"github.com/gorilla/mux"
	"encoding/json"
)

func StartServer() {
	fmt.Println("[INFO] Initializing database")
	db.Init()

	// Démarre le scan périodique
	go StartPeriodicScan()

	fmt.Println("[INFO] Starting IoT Dashboard Server on :8080")
	r := mux.NewRouter()

	r.HandleFunc("/api/devices/{id}/action", ControlDeviceHandler).Methods("POST")
	r.HandleFunc("/api/devices/{id}", DeleteDeviceHandler).Methods("DELETE")
	r.HandleFunc("/api/devices/{id}", UpdateDeviceHandler).Methods("PUT")
	r.HandleFunc("/api/scan", ScanHandler)
	r.HandleFunc("/api/devices", GetDevicesHandler)
	r.HandleFunc("/api/add-mock-device", AddMockDeviceHandler)

	http.Handle("/", r)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler pour scanner le réseau
func ScanHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := scanner.ScanNetwork()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, d := range devices {
		db.InsertDevice(d)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// Handler pour obtenir tous les devices
func GetDevicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	devices := db.GetAllDevices()
	json.NewEncoder(w).Encode(devices)
}

// Handler pour contrôler un device
func ControlDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	device := db.GetDeviceByID(id)
	if device.ID == 0 {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	var payload struct {
		Action string `json:"action"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	result, err := PerformAction(device, payload.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(result))
}

// Handler suppression device
func DeleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	_, err = db.Database().Exec("DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Erreur suppression", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Handler modification device
func UpdateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	var payload struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	_, err = db.Database().Exec("UPDATE devices SET name = ?, desc = ? WHERE id = ?", payload.Name, payload.Desc, id)
	if err != nil {
		http.Error(w, "Erreur update", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Scan périodique en tâche de fond
func StartPeriodicScan() {
	go func() {
		for {
			devices, err := scanner.ScanNetwork()
			if err == nil {
				var ips []string
				for _, d := range devices {
					existing := db.GetDeviceByIP(d.IP)
					if existing.IsMock {
						continue
					}
					d.Online = true
					db.InsertDevice(d)
					ips = append(ips, d.IP)
				}
				all := db.GetAllDevices()
				for _, dev := range all {
					if dev.IsMock {
						db.SetDeviceOnlineStatus(dev.IP, true)
						ips = append(ips, dev.IP)
					} else if !contains(ips, dev.IP) {
						db.SetDeviceOnlineStatus(dev.IP, false)
					}
				}
				db.RemoveDevicesNotInList(ips)
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
