package model

// ClientInfo represents client network metadata
type ClientInfo struct {
	IP      string `json:"ip"`
	ISP     string `json:"isp"`
	Country string `json:"country"`
	Lat     string `json:"lat"`
	Lon     string `json:"lon"`
}
