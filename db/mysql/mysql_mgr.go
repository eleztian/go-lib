package mysql

import (
	"database/sql"
	"log"
	"time"
)

type MysqlMgr struct {
	session *sql.DB
	timeout time.Duration
	dbcfg   *MysqlConnInfo
}

func NewMysqlMgr(dbcfg *MysqlConnInfo, timeout time.Duration) *MysqlMgr {
	return &MysqlMgr{
		session: nil,
		timeout: timeout,
		dbcfg:   dbcfg,
	}
}

// Start 连接mysql
func (mgr *MysqlMgr) Start() error {
	var err error
	mgr.session, err = sql.Open("mysql", mgr.dbcfg.String())
	if err != nil {
		return err
	}
	if err = mgr.Ping(); err != nil {
		return err
	}
	log.Println("mysql", "mysql connect success.")
	return nil
}

func (mgr *MysqlMgr) Close() {
	if mgr.session != nil {
		mgr.session.Close()
		mgr.session = nil
		log.Println("mysql", "mysql disconnect")
	}
}

func (mgr *MysqlMgr) Ping() error {
	if mgr.session != nil {
		return mgr.session.Ping()
	}
	return MYSQL_SESSION_NIL_ERR
}

func (mgr *MysqlMgr) GetSession() *sql.DB {
	return mgr.session
}

// GetNamedRows 解析查询返回的数据
func (mgr *MysqlMgr) GetNamedRows(query interface{}) ([]map[string]interface{}, error) {
	return getNamedRows(query)
}

func getNamedRows(query interface{}) ([]map[string]interface{}, error) {
	row, ok := query.(*sql.Rows)
	var results []map[string]interface{}
	if ok == false {
		return results, MYSQL_GET_NAMED_RESULT_ERROR
	}
	columnTypes, err := row.ColumnTypes()
	if err != nil {
		log.Printf("mysql columns error:%v\n", err)
		return nil, err
	}

	values := make([]interface{}, len(columnTypes))

	for c := true; c || row.NextResultSet(); c = false {

		//maybe mgr way is better
		//for k, c := range columnTypes {
		//	scans[k] = reflect.New(c.ScanType()).Interface()
		//}

		for k, c := range columnTypes {
			switch c.DatabaseTypeName() {
			case "TINYINT":
				values[k] = new(int8)
			case "INT":
				values[k] = new(int32)
			case "BIGINT":
				values[k] = new(int64)
			case "DOUBLE", "FLOAT":
				values[k] = new(float64)
			case "VARCHAR":
				values[k] = new(string)
			case "BLOB":
				values[k] = new([]byte)
			default:
				values[k] = new([]byte)
			}
		}

		for row.Next() {
			if err = row.Scan(values...); err != nil {
				log.Printf("mysql Scan error:%v\n", err)
				return nil, err
			}

			result := make(map[string]interface{})
			for k, v := range values {
				key := columnTypes[k]
				switch s := v.(type) {
				case *int8:
					result[key.Name()] = *s
				case *int64:
					result[key.Name()] = *s
				case *float64:
					result[key.Name()] = *s
				case *string:
					result[key.Name()] = *s
				case *[]byte:
					result[key.Name()] = *s
				}
			}
			results = append(results, result)
		}
	}
	return results, nil
}
