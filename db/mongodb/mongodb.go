package mongodb

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type M = bson.M

const (
	Strong    = 1
	Monotonic = 2
)

type MongoDBMgr struct {
	session *mgo.Session
	timeout time.Duration
	dbcfg   *MongoDBInfo
}

func NewMongoDBMgr(dbcfg *MongoDBInfo, timeout time.Duration) *MongoDBMgr {
	return &MongoDBMgr{
		session: nil,
		timeout: timeout,
		dbcfg:   dbcfg,
	}
}

func (mgr *MongoDBMgr) Start() error {
	var err error
	mgr.session, err = mgo.DialWithTimeout(mgr.dbcfg.String(), mgr.timeout)
	if err != nil {
		return err
	}
	if err = mgr.Ping(); err != nil {
		return err
	}
	mgr.session.SetMode(Monotonic, true)
	fmt.Println("mongodb", "MongoDB Connect success")
	return nil
}

func (mgr *MongoDBMgr) Close() {
	if mgr.session != nil {
		mgr.session.DB("").Logout()
		mgr.session.Close()
		mgr.session = nil
		fmt.Println("mongodb", "Disconnect mongodb url: ", mgr.dbcfg.String())
	}
}

func (mgr *MongoDBMgr) Ping() error {
	if mgr.session != nil {
		return mgr.session.Ping()
	}
	return MONGODB_SESSION_NIL_ERR
}

func (mgr *MongoDBMgr) RefreshSession() {
	mgr.session.Refresh()

}

func (mgr *MongoDBMgr) GetDbSession() *mgo.Session {
	return mgr.session
}

func (mgr *MongoDBMgr) SetMode(mode int, refresh bool) {
	status := mgo.Monotonic
	if mode == Strong {
		status = mgo.Strong
	} else {
		status = mgo.Monotonic
	}

	mgr.session.SetMode(status, refresh)
}

func (mgr *MongoDBMgr) Insert(collection string, doc interface{}) error {
	if mgr.session == nil {
		return MONGODB_SESSION_NIL_ERR
	}

	dbSession := mgr.session.Copy()
	defer dbSession.Close()

	c := dbSession.DB("").C(collection)

	return c.Insert(doc)
}

func (mgr *MongoDBMgr) Update(collection string, cond interface{}, change interface{}) error {
	if mgr.session == nil {
		return MONGODB_SESSION_NIL_ERR
	}

	dbSession := mgr.session.Copy()
	defer dbSession.Close()

	c := dbSession.DB("").C(collection)

	return c.Update(cond, bson.M{"$set": change})
}

func (mgr *MongoDBMgr) UpdateInsert(collection string, cond interface{}, doc interface{}) error {
	if mgr.session == nil {
		return MONGODB_SESSION_NIL_ERR
	}

	dbSession := mgr.session.Copy()
	defer dbSession.Close()

	c := dbSession.DB("").C(collection)
	_, err := c.Upsert(cond, bson.M{"$set": doc})
	if err != nil {
		fmt.Printf("mongodb UpdateInsert failed collection is:%s. cond is:%v\n", collection, cond)
	}

	return err
}

func (mgr *MongoDBMgr) RemoveOne(collection string, condName string, condValue int64) error {
	if mgr.session == nil {
		return MONGODB_SESSION_NIL_ERR
	}

	dbSession := mgr.session.Copy()
	defer dbSession.Close()

	c := dbSession.DB("").C(collection)
	err := c.Remove(bson.M{condName: condValue})
	if err != nil && err != mgo.ErrNotFound {
		fmt.Printf("mongodb remove failed from collection:%s. name:%s-value:%d\n", collection, condName, condValue)
	}

	return err

}

func (mgr *MongoDBMgr) RemoveOneByCond(collection string, cond interface{}) error {
	if mgr.session == nil {
		return MONGODB_SESSION_NIL_ERR
	}

	dbSession := mgr.session.Copy()
	defer dbSession.Close()

	c := dbSession.DB("").C(collection)
	err := c.Remove(cond)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Printf("mongodb remove failed from collection:%s. cond :%v, err: %v.", collection, cond, err)
	}

	return err

}

func (mgr *MongoDBMgr) RemoveAll(collection string, cond interface{}) error {
	if mgr.session == nil {
		return MONGODB_SESSION_NIL_ERR
	}

	dbSession := mgr.session.Copy()
	defer dbSession.Close()

	c := dbSession.DB("").C(collection)
	change, err := c.RemoveAll(cond)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Printf("mongodb MongoDBMgr RemoveAll failed : %s, %v", collection, cond)
		return err
	}
	fmt.Printf("mongodb MongoDBMgr RemoveAll: %v, %v", change.Updated, change.Removed)
	return nil
}

func (mgr *MongoDBMgr) DBQuery(collection string, cond interface{}, result *[]map[string]interface{}) error {
	if mgr.session == nil {
		return MONGODB_SESSION_NIL_ERR
	}

	dbSession := mgr.session.Copy()
	defer dbSession.Close()

	c := dbSession.DB("").C(collection)
	q := c.Find(cond)

	if nil == q {
		return MONGODB_DBFINDALL_ERR
	}

	q.All(result)
	return nil
}

func (mgr *MongoDBMgr) DBQueryOneResult(collection string, cond interface{}, result map[string]interface{}) error {
	if mgr.session == nil {
		return MONGODB_SESSION_NIL_ERR
	}

	dbSession := mgr.session.Copy()
	defer dbSession.Close()

	c := dbSession.DB("").C(collection)
	q := c.Find(cond)

	if nil == q {
		return MONGODB_DBFINDALL_ERR
	}

	q.One(result)
	return nil
}
