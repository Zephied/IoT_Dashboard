package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"iot-dashboard/backend/models"
)

type ActionConfig struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func PerformAction(device models.Device, action string) (string, error) {
	// Traitement spécial pour le mock caméra locale
	if device.IP == "127.0.0.2" && action == "voir la caméra" {
		// Retourne une fausse URL vidéo (à servir statiquement côté frontend ou backend)
		return "/mock-stream.mp4", nil
	}

	var actions map[string]map[string]interface{}
	err := json.Unmarshal([]byte(device.Actions), &actions)
	if err != nil {
		return "", fmt.Errorf("invalid actions format")
	}

	act, exists := actions[action]
	if !exists {
		return "", fmt.Errorf("action not found")
	}

	// Substitution {{ip}} dans l’URL
	url := act["url"].(string)
	url = strings.ReplaceAll(url, "{{ip}}", device.IP)

	// Méthode HTTP
	method := act["method"].(string)

	// Corps
	bodyStr := ""
	if b, ok := act["body"].(string); ok {
		bodyStr = b
	}
	body := bytes.NewBufferString(bodyStr)

	// Headers
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}
	if headers, ok := act["headers"].(map[string]interface{}); ok {
		for k, v := range headers {
			req.Header.Set(k, fmt.Sprint(v))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return string(respBody), nil
}

func renderTemplate(tmplStr string, data map[string]string) (string, error) {
	tmpl, err := template.New("tpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
