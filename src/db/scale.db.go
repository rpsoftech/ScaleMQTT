package db

import (
	"encoding/json"
	"rpsoftech/scaleMQTT/src/systypes"
)

const preFixKeyForUsernameAndPassword = "scaleusernamepass/"

// const preFixKeyForUsernameAndPassword = "scaleconfig/"
const preFixKeyForScaleConfig = "scaleconfig/"

func (DBObject *DbClass) AddScaleUserNamePassword(username string, password []byte) error {
	return DBClassObject.connection.Put([]byte(preFixKeyForUsernameAndPassword+username), password)
}

func (DBObject *DbClass) AddScaleConfigData(username string, config systypes.ScaleConfigData) (bool, error) {
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

func (DBObject *DbClass) GetScaleConfigData(username string) (config systypes.ScaleConfigData, err error) {
	stringConfig, err := DBClassObject.connection.Get([]byte(preFixKeyForUsernameAndPassword + username))
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
func (DBObject *DbClass) GetPasswordForScale(username string) ([]byte, error) {
	return DBClassObject.connection.Get([]byte(preFixKeyForUsernameAndPassword + username))
}
