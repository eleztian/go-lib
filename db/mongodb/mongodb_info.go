package mongodb

import (
	"errors"
	"fmt"
)

var (
	MONGODB_SESSION_NIL_ERR = errors.New("MongoDBMgr session nil.")
	MONGODB_NOTFOUND_ERR    = errors.New("not found!")
	MONGODB_DBFINDALL_ERR   = errors.New("MongoDBMgr found error")
)

type MongoDBInfo struct {
	DbHost string
	DbPort int
	DbName string
	DbUser string
	DbPass string
}

func NewMongoDBInfo(host string, port int, name, user, pass string) *MongoDBInfo {
	return &MongoDBInfo{
		DbHost: host,
		DbPort: port,
		DbName: name,
		DbUser: user,
		DbPass: pass,
	}
}

func (mi *MongoDBInfo) String() string {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		mi.DbUser, mi.DbPass, mi.DbHost, mi.DbPort, mi.DbName)
	if mi.DbUser == "" || mi.DbPass == "" {
		url = fmt.Sprintf("mongodb://%s:%d/%s", mi.DbHost, mi.DbPort, mi.DbName)
	}
	return url
}
