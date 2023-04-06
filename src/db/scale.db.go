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

func (DBObject *DbClass) AddScaleUserNamePassword(username string, password []byte) error {
	return DBClassObject.connection.Put([]byte(preFixKeyForUsernameAndPassword+username), password)
}

func (DBObject *DbClass) AddScaleConfigData(username string, config ScaleConfigData) (bool, error) {
	stringConfig, err := json.Marshal(config)
	if err != nil {
		return false, err
	}
	err = DBClassObject.connection.Put([]byte(preFixKeyForScaleConfig+username), stringConfig)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (DBObject *DbClass) GetScaleConfigData(username string) (config ScaleConfigData, err error) {
	stringConfig, err := DBClassObject.connection.Get([]byte(preFixKeyForUsernameAndPassword + username))
	if err != nil {
		return
	}
	err = json.Unmarshal(stringConfig, &config)
	return
}
func (DBObject *DbClass) GetPasswordForScale(username string) ([]byte, error) {
	return DBClassObject.connection.Get([]byte(preFixKeyForUsernameAndPassword + username))
}
