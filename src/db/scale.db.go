package db

import (
	"encoding/json"
)

const preFixKeyForUsernameAndPassword = "scaleusernamepass/"
const preFixKeyForScaleConfig = "scaleconfig/"

type ScaleConfigData struct {
	WifiSsid           string `json:"wifi_ssid" validator:"requried" binding:"requried"`
	WifiPassword       string `json:"wifi_password" validator:"requried" binding:"requried"`
	MqttBrokerURI      string `json:"mqtt_broker_uri" validator:"requried" binding:"requried"`
	MqttUsername       string `json:"mqtt_username" validator:"requried" binding:"requried"`
	MqttPassword       string `json:"mqtt_password" validator:"requried" binding:"requried"`
	MqttPublishTopic   string `json:"mqtt_publish_topic" validator:"requried" binding:"requried"`
	MqttSubscribeTopic string `json:"mqtt_subscribe_topic" validator:"requried" binding:"requried"`
	MqttBrokerPort     int    `json:"mqtt_broker_port" validator:"requried" binding:"requried"`
	BaudRate           int    `json:"baud_rate" validator:"requried" binding:"requried"`
	LogEnable          bool   `json:"log_enable" validator:"requried" binding:"requried"`
}

func AddScaleUserNamePassword(username string, password []byte) error {
	return DbConnection.Put([]byte(preFixKeyForUsernameAndPassword+username), password)
}

func AddScaleConfigData(username string, config ScaleConfigData) (bool, error) {
	stringConfig, err := json.Marshal(config)
	if err != nil {
		return false, err
	}
	err = DbConnection.Put([]byte(preFixKeyForScaleConfig+username), stringConfig)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetScaleConfigData(username string) (config ScaleConfigData, err error) {
	stringConfig, err := DbConnection.Get([]byte(preFixKeyForUsernameAndPassword + username))
	if err != nil {
		return
	}
	err = json.Unmarshal(stringConfig, &config)
	return
}
func GetPasswordForScale(username string) ([]byte, error) {
	return DbConnection.Get([]byte(preFixKeyForUsernameAndPassword + username))
}
