package db

import (
	"path/filepath"
	"rpsoftech/scaleMQTT/src/global"
	"strings"

	"git.mills.io/prologic/bitcask"
)

var DbConnection *bitcask.Bitcask

func init() {
	db, _ := bitcask.Open(filepath.Join(global.GetCuurentPath(), "dbcollection"))
	DbConnection = db
}

func TakeBackup() map[string]map[string]string {
	data := make(map[string]map[string]string)
	DbConnection.Fold(func(key []byte) error {

		val, _ := DbConnection.Get(key)

		indexes := strings.Split(string(key), "/")
		if _, ok := data[indexes[0]]; !ok {
			data[indexes[0]] = make(map[string]string)
		}
		data[indexes[0]][string(indexes[1])] = string(val)
		return nil
	})
	return data
}
