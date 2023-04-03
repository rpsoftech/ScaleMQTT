package global

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

var JWTKEY []byte
var Logger *zerolog.Logger

func GetCuurentPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
