package db

import (
	"github.com/eleztian/go-lib/db/mysql"
	"testing"
	"time"
)

func TestMysql(t *testing.T) {
	mgr := mysql.NewMysqlMgr(
		mysql.NewMysqlConnInfo("127.0.0.1", 3306, "name", "****", "****"),
		10*time.Second)
	err := mgr.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer mgr.Close()

	rows, err := mgr.GetSession().Query("select * from product_info")
	if err != nil {
		t.Error(err)
		return
	}
	rows.Close()
	m, err := mgr.GetNamedRows(rows)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range m {
		t.Log(v)
	}

}
