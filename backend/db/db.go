package db

import (
	"database/sql"
	"log"

	"iot-dashboard/backend/models"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func Init() {
	var err error
	database, err = sql.Open("sqlite3", "./iot.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS devices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT UNIQUE,
		mac TEXT UNIQUE,
		name TEXT,
		desc TEXT,
		type TEXT,
		actions TEXT,
		online INTEGER DEFAULT 1,
		is_mock INTEGER DEFAULT 0
	);`
	_, err = database.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Erreur création de la table devices:", err)
	}
}

func Database() *sql.DB {
	return database
}

func InsertDevice(d models.Device) {
	stmt, err := database.Prepare("INSERT OR REPLACE INTO devices(ip, mac, name, desc, type, actions, online, is_mock) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Erreur prepare insert:", err)
		return
	}
	_, err = stmt.Exec(d.IP, d.MAC, d.Name, d.Desc, d.Type, d.Actions, boolToInt(d.Online), boolToInt(d.IsMock))
	if err != nil {
		log.Println("Erreur insert device:", err)
	}
}

func GetAllDevices() []models.Device {
	rows, err := database.Query("SELECT id, ip, mac, name, desc, type, actions, online, is_mock FROM devices")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var d models.Device
		var onlineInt, isMockInt int
		err := rows.Scan(&d.ID, &d.IP, &d.MAC, &d.Name, &d.Desc, &d.Type, &d.Actions, &onlineInt, &isMockInt)
		if err != nil {
			panic(err)
		}
		d.Online = onlineInt == 1
		d.IsMock = isMockInt == 1
		devices = append(devices, d)
	}
	return devices
}

func GetDeviceByID(id int) models.Device {
	var device models.Device
	var onlineInt, isMockInt int
	row := database.QueryRow("SELECT id, ip, mac, name, desc, type, actions, online, is_mock FROM devices WHERE id = ?", id)
	err := row.Scan(&device.ID, &device.IP, &device.MAC, &device.Name, &device.Desc, &device.Type, &device.Actions, &onlineInt, &isMockInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Device{}
		}
		log.Println("Erreur GetDeviceByID:", err)
		return models.Device{}
	}
	device.Online = onlineInt == 1
	device.IsMock = isMockInt == 1
	return device
}

func GetDeviceByIP(ip string) models.Device {
	var device models.Device
	var onlineInt, isMockInt int
	row := database.QueryRow("SELECT id, ip, mac, name, desc, type, actions, online, is_mock FROM devices WHERE ip = ?", ip)
	err := row.Scan(&device.ID, &device.IP, &device.MAC, &device.Name, &device.Desc, &device.Type, &device.Actions, &onlineInt, &isMockInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Device{}
		}
		log.Println("Erreur GetDeviceByIP:", err)
		return models.Device{}
	}
	device.Online = onlineInt == 1
	device.IsMock = isMockInt == 1
	return device
}

func SetDeviceOnlineStatus(ip string, online bool) {
	stmt, err := database.Prepare("UPDATE devices SET online = ? WHERE ip = ?")
	if err != nil {
		log.Println("Erreur prepare update online:", err)
		return
	}
	_, err = stmt.Exec(boolToInt(online), ip)
	if err != nil {
		log.Println("Erreur update online:", err)
	}
}

func RemoveDevicesNotInList(ips []string) {
	if len(ips) == 0 {
		return
	}
	// On ne supprime que les devices non mocks
	query := "DELETE FROM devices WHERE ip NOT IN (" + placeholders(len(ips)) + ") AND is_mock = 0"
	args := make([]interface{}, len(ips))
	for i, ip := range ips {
		args[i] = ip
	}
	_, err := database.Exec(query, args...)
	if err != nil {
		log.Println("Erreur suppression devices non détectés:", err)
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func placeholders(n int) string {
	if n == 0 {
		return "''"
	}
	ph := "?"
	for i := 1; i < n; i++ {
		ph += ",?"
	}
	return ph
}
