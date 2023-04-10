package db

import (
	"encoding/json"
	"rpsoftech/scaleMQTT/src/systypes"
)

const preFixKeyForUsernameAndPassword = "scaleusernamepass/"

// const preFixKeyForUsernameAndPassword = "scaleconfig/"
const preFixKeyForScaleConfig = "scaleconfig/"

func AddScaleUserNamePassword(username string, password []byte) error {
	return DbConnection.Put([]byte(preFixKeyForUsernameAndPassword+username), password)
}

func AddScaleConfigData(deviceUniqueID string, config systypes.ScaleConfigData) (bool, error) {
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

func GetScaleConfigData(deviceUniqueID string) (config systypes.ScaleConfigData, err error) {
	stringConfig, err := DbConnection.Get([]byte(preFixKeyForScaleConfig + deviceUniqueID))
	if err != nil {
		return
	}
	err = json.Unmarshal(stringConfig, &config)
	if err != nil {
		return
	}
	_, err = config.Validate()
	return
}
func GetPasswordForScale(username string) ([]byte, error) {
	return DbConnection.Get([]byte(preFixKeyForUsernameAndPassword + username))
}
