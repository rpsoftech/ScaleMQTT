package global

import (
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type MQTTConnectionMeta struct {
	Connected  bool    `json:"connected"`
	UserName   string  `json:"name"`
	LocationID string  `json:"locationName"`
	Weight     float64 `json:"weight"`
	Count      int     `json:"-"`
}

var JWTKEY []byte
var Logger *zerolog.Logger
var Validator = validator.New()

var MQTTConnectionStatusMap = make(map[string]*MQTTConnectionMeta)

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
