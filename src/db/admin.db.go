package db

func (DBObject *DbClass) AddNewAdmin(username []byte, password []byte) error {
	return DBObject.connection.Put(username, password)
}
func (DBObject *DbClass) GetAdmin(username []byte) ([]byte, error) {
	return DBObject.connection.Get(username)
}
