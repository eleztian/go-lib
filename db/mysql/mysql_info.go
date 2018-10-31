package mysql

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

var (
	MYSQL_SESSION_NIL_ERR        = errors.New("MysqlMgr session nil.")
	MYSQL_GET_NAMED_RESULT_ERROR = errors.New("Mysql GetResult Type Error")
)

type MysqlConnInfo struct {
	*mysql.Config
}

func NewMysqlConnInfo(host string, port int, name, user, pass string) *MysqlConnInfo {
	cfg := mysql.NewConfig()
	cfg.User = user
	cfg.Passwd = pass
	cfg.Addr = fmt.Sprintf("%s:%d", host, port)
	cfg.Net = "tcp"
	cfg.DBName = name
	cfg.MultiStatements = true
	cfg.InterpolateParams = true
	return &MysqlConnInfo{cfg}
}

func (mi *MysqlConnInfo) String() string {
	return mi.FormatDSN()
}
