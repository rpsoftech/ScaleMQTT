package db

import (
	"path/filepath"
	"rpsoftech/scaleMQTT/src/global"
	"strings"

	"git.mills.io/prologic/bitcask"
)

type DbClass struct {
	connection *bitcask.Bitcask
}

var DBClassObject *DbClass

func (DBObject *DbClass) CloseConnection() {
	DBObject.connection.Close()
}
func (DBObject *DbClass) SetConnection(con *bitcask.Bitcask) {
	DBObject.connection = con
}
func init() {
	db, _ := bitcask.Open(filepath.Join(global.GetCuurentPath(), "dbcollection"))
	DBClassObject = &DbClass{
		connection: db,
	}
}

func (DBObject *DbClass) TakeBackup() map[string]map[string]string {
	data := make(map[string]map[string]string)
	DBObject.connection.Fold(func(key []byte) error {

		val, _ := DBObject.connection.Get(key)

		indexes := strings.Split(string(key), "/")
		if _, ok := data[indexes[0]]; !ok {
			data[indexes[0]] = make(map[string]string)
		}
		data[indexes[0]][string(indexes[1])] = string(val)
		return nil
	})
	return data
}
