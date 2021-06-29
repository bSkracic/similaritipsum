package db

import (
	"sync"

	"github.com/bSkracic/similaritipsum/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Conn struct {
	*gorm.DB
}

var (
	db   *Conn
	once sync.Once
)

func GetConnection() *Conn {
	once.Do(func() {
		dbCfg := config.GetFromEnv()
		tempDB, err := gorm.Open(postgres.Open(dbCfg.ConnString()), &gorm.Config{})
		if err != nil {
			panic("Database connection failed")
		}
		db = &Conn{tempDB}
	})
	return db
}

func (db *Conn) WordEntries() *gorm.DB {
	return db.Table("word_entries")
}
