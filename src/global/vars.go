package global

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"rpsoftech/scaleMQTT/src/systypes"

	"github.com/joho/godotenv"
	"github.com/mochi-co/mqtt/v2"
	"github.com/rs/zerolog"
)

var JWTKEY []byte
var Logger *zerolog.Logger
var MQTTserver *mqtt.Server

var MQTTConnectionStatusMap = make(map[string]*systypes.MQTTConnectionMeta)
var MQTTConnectionWithUidStatusMap = make(map[string]*systypes.MQTTConnectionMeta)

const DefaultMQTTDevicePublishTopicSuffix = "/" + systypes.DefaultMQTTDevicePublishTopicSuffix
const DefaultMQTTDeviceSubscribeTopicSuffix = "/" + systypes.DefaultMQTTDeviceSubscribeTopicSuffix
const RootDefaultDeviceId = "thisisdefaultidforiotdevice"

func init() {
	LoadEnv()
}
func LoadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		// path/to/whatever exists
		godotenv.Load(".env")
	} else {
		godotenv.Load("./../.env")
	}
	defaultValue := make([]byte, 128)

	_, err := rand.Read(defaultValue)
	if err != nil {
		defaultValue = []byte("thisisjustdefaultvalue")
	}
	defaultValueString := hex.EncodeToString(defaultValue)
	envJWTKeyValue := os.Getenv("JWTKEY")
	if envJWTKeyValue == "" {
		envJWTKeyValue = defaultValueString
	}
	JWTKEY = []byte(envJWTKeyValue)
}

func IsConnected(username string) (ok bool) {
	if val, has := MQTTConnectionStatusMap[username]; has {
		ok = val.Connected
	} else {
		ok = false
	}
	return
}

func GetCuurentPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
