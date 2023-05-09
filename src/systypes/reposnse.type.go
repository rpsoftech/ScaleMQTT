package systypes

type BaseResponseFormat struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type MQTTConnectionMeta struct {
	Config     ScaleConfigData `json:"-"`
	Connected  bool            `json:"connected"`
	UserName   string          `json:"name"`
	LocationID string          `json:"locationName"`
	Weight     float64         `json:"weight"`
	RawWeight  string          `json:"raw_info"`
	Count      int             `json:"-"`
}
