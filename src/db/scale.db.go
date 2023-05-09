package db

import (
	"encoding/json"
	"rpsoftech/scaleMQTT/src/systypes"
)

const preFixKeyForUsernameAndPassword = "scaleusernamepass/"

const preFixForOldToNewDeviceId = "scaleDeviceIdMap/"
const preFixKeyForScaleConfig = "scaleconfig/"

func (DBObject *DbClass) GetScaleConfigData(devID string) (config systypes.ScaleConfigData, err error) {
	stringConfig, err := DBClassObject.connection.Get([]byte(preFixKeyForScaleConfig + devID))
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
func (DBObject *DbClass) SetScaleConfigData(devID string, config systypes.ScaleConfigData) (err error) {
	var stringConfig []byte
	stringConfig, err = config.JSON()
	if err != nil {
		return
	}
	err = DBClassObject.connection.Put([]byte(preFixKeyForScaleConfig+devID), stringConfig)
	if err != nil {
		return
	}
	return
}

//	func (DBObject *DbClass) GetPasswordForScale(username string) ([]byte, error) {
//		return DBClassObject.connection.Get([]byte(preFixKeyForUsernameAndPassword + username))
//	}
func (DBObject *DbClass) SetChanegedDeviceId(OldDevID string, NewDeviceId string) error {
	return DBClassObject.connection.Put([]byte(preFixForOldToNewDeviceId+OldDevID), []byte(NewDeviceId))
}

func (DBObject *DbClass) GetChanegedDeviceId(username string) ([]byte, error) {
	return DBClassObject.connection.Get([]byte(preFixKeyForUsernameAndPassword + username))
}
