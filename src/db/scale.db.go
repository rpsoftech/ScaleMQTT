package db

import (
	"encoding/json"
	"rpsoftech/scaleMQTT/src/global"
)

const preFixKeyForUsernameAndPassword = "scaleusernamepass/"

// const preFixKeyForUsernameAndPassword = "scaleconfig/"
const preFixKeyForScaleConfig = "scaleconfig/"

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
	DivideMultiply   string `json:"divide_multiply" validate:"required" binding:"required"`
	DivideMultiplyBy int    `json:"divide_multiply_by" validate:"required,int" binding:"required"`
}

func AddScaleUserNamePassword(username string, password []byte) error {
	return DbConnection.Put([]byte(preFixKeyForUsernameAndPassword+username), password)
}

func AddScaleConfigData(deviceUniqueID string, config ScaleConfigData) (bool, error) {
	stringConfig, err := json.Marshal(config)
	if err != nil {
		return false, err
	}
	err = DbConnection.Put([]byte(preFixKeyForScaleConfig+deviceUniqueID), stringConfig)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetScaleConfigData(deviceUniqueID string) (config ScaleConfigData, err error) {
	stringConfig, err := DbConnection.Get([]byte(preFixKeyForScaleConfig + deviceUniqueID))
	if err != nil {
		return
	}
	err = json.Unmarshal(stringConfig, &config)
	if err != nil {
		return
	}
	_, err = config.validate()
	return
}
func GetPasswordForScale(username string) ([]byte, error) {
	return DbConnection.Get([]byte(preFixKeyForUsernameAndPassword + username))
}

func (data *ScaleConfigData) validate() (valid bool, err error) {
	err = global.Validator.Struct(data)
	if err == nil {
		valid = true
	}
	return
}
