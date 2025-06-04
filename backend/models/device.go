package models

type Device struct {
	ID      int    `json:"id"`
	IP      string `json:"ip"`
	MAC     string `json:"mac"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Type    string `json:"type"`
	Actions string `json:"actions"`
	Online  bool   `json:"online"`
	IsMock  bool   `json:"isMock"`
}
