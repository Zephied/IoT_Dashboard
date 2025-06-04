package scanner

import (
	"fmt"
	"iot-dashboard/backend/models"

	"github.com/Ullaakut/nmap"
)

func ScanNetwork() ([]models.Device, error) {
	scanner, err := nmap.NewScanner(
		nmap.WithTargets("192.168.1.0/24"),
		nmap.WithPingScan(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create scanner: %w", err)
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}
	if warnings != nil {
		fmt.Printf("[WARN] %v\n", warnings)
	}

	var devices []models.Device
	for _, host := range result.Hosts {
		if len(host.Addresses) == 0 {
			continue
		}
		var mac string
		for _, addr := range host.Addresses {
			if addr.AddrType == "mac" {
				mac = addr.Addr
				break
			}
		}
		name := ""
		if len(host.Hostnames) > 0 {
			name = host.Hostnames[0].Name
		}
		device := models.Device{
			IP:   host.Addresses[0].String(),
			MAC:  mac,
			Name: name,
		}
		devices = append(devices, device)
	}

	return devices, nil
}
