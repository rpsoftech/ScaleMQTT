package systypes

import (
	"bytes"
	"encoding/json"
	"errors"
	"rpsoftech/scaleMQTT/src/validator"
)

type DivideMultiplyString string

const DefaultMQTTDevicePublishTopicSuffix = "WeighingScale/DeviceConfig/up"
const DefaultMQTTDeviceSubscribeTopicSuffix = "WeighingScale/DeviceConfig/dw"
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

type MqttBrokerConfig struct {
	MqttBrokerURI          string `json:"broker_uri" validate:"required" binding:"required"`
	MqttUsername           string `json:"mqtt_username" validate:"required" binding:"required"`
	MqttPassword           string `json:"mqtt_password" validate:"required" binding:"required"`
	MqttPublishTopic       string `json:"pub_topic" validate:"required" binding:"required"`
	MqttSerialPublishTopic string `json:"serial_pub_topic" validate:"required" binding:"required"`
	MqttSubscribeTopic     string `json:"sub_topic" validate:"required" binding:"required"`
	MqttBrokerPort         int    `json:"broker_port" validate:"required,port" binding:"required"`
}
type DevcfgForMqtt struct {
	MqttBrokerConfig
	DevID        string `json:"dev_id" validate:"required" binding:"required"`
	WifiSsid     string `json:"wifi_ssid" validate:"required" binding:"required"`
	WifiPassword string `json:"wifi_password" validate:"required" binding:"required"`
	EndChar      string `json:"end_char" validate:"required" binding:"required"`
	BaudRate     int    `json:"baud_rate" validate:"required" binding:"required"`
	LogEnable    bool   `json:"log_enable" validate:"required" binding:"required"`
}

type ScaleConfigData struct {
	DevcfgForMqtt
	DivideMultiply   DivideMultiplyString `json:"divide_multiply" validate:"required" binding:"required"`
	DivideMultiplyBy int                  `json:"divide_multiply_by,omitempty"`
	NegativeChar     string               `json:"negative_char"`
}
type ChangeDeviceIdData struct {
	DevID    string `json:"dev_id" validate:"required" binding:"required"`
	OldDevID string `json:"old_dev_id" validate:"required" binding:"required"`
}

func (data *ScaleConfigData) Validate() (valid bool, err error) {
	data.MqttSubscribeTopic = DefaultMQTTDeviceSubscribeTopicSuffix
	data.MqttPublishTopic = DefaultMQTTDevicePublishTopicSuffix
	if data.NegativeChar == "" {
		data.NegativeChar = "\\f"
	}
	err = validator.Validator.Struct(data)
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

func (t *ScaleConfigData) JSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
