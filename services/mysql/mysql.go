package mysql

import (
	"context"
	"database/sql"
	"sync"

	"gopkg.in/gorp.v2"
	"test.com/mine/services/initializer"

	_ "github.com/go-sql-driver/mysql"
)

var (
	all []initializer.Simple
)

type Manager struct {
}

func (m *Manager) GetConn() *gorp.DbMap {
	return dbMap
}

type initMysql struct {
}

var (
	once  = sync.Once{}
	dbMap *gorp.DbMap
)

func (initMysql) Initial(context.Context) {
	once.Do(func() {
		db, err := sql.Open("mysql", "root:my-secret-pw@tcp(127.0.0.1:3306)/food?charset=utf8&parseTime=true")
		if err != nil {
			panic(err)
		}
		dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
		err = dbMap.Db.Ping()
		if err != nil {
			panic(err)
		}

		for i := range all {
			all[i].Initialize()
		}
	})
}

// Register a new initMysql module
func Register(m ...initializer.Simple) {
	all = append(all, m...)
}

func init() {
	initializer.Register(initMysql{}, 0)
}
