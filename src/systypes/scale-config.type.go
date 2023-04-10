package systypes

import "errors"

type DivideMultiplyString string

const (
	Divide DivideMultiplyString = "/"
	Multi  DivideMultiplyString = "*"
)

func (s DivideMultiplyString) String() string {
	switch s {
	case Divide:
		return "/"
	case Multi:
		return "*"
	}
	return "unknown"
}

type DevcfgForMqtt struct {
	DevID              string `json:"dev_id" validate:"required" binding:"required"`
	WifiSsid           string `json:"wifi_ssid" validate:"required" binding:"required"`
	WifiPassword       string `json:"wifi_password" validate:"required" binding:"required"`
	MqttBrokerURI      string `json:"mqtt_broker_uri" validate:"required" binding:"required"`
	MqttUsername       string `json:"mqtt_username" validate:"required" binding:"required"`
	MqttPassword       string `json:"mqtt_password" validate:"required" binding:"required"`
	MqttPublishTopic   string `json:"mqtt_publish_topic" validate:"required" binding:"required"`
	MqttSubscribeTopic string `json:"mqtt_subscribe_topic" validate:"required" binding:"required"`
	MqttBrokerPort     int    `json:"mqtt_broker_port" validate:"required" binding:"required"`
	EndChar            string `json:"end_char" validate:"required" binding:"required"`
	BaudRate           int    `json:"baud_rate" validate:"required" binding:"required"`
	LogEnable          bool   `json:"log_enable" validate:"required" binding:"required"`
}

type ScaleConfigData struct {
	DevcfgForMqtt
	DivideMultiply   DivideMultiplyString `json:"divide_multiply" validate:"required" binding:"required"`
	DivideMultiplyBy int                  `json:"divide_multiply_by,omitempty"`
}

func (data *ScaleConfigData) Validate() (valid bool, err error) {

	err = Validator.Struct(data)
	if err == nil {
		valid = true
	} else {
		valid = false
		return
	}
	if data.DivideMultiply == Divide {
		if data.DivideMultiplyBy == 0 {
			err = errors.New("can not divide by 0")
			return
		}
	}
	return
}
