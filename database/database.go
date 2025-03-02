package database

import (
	"fmt"
	"net/url"

	"github.com/BlenDMinh/dutgrad-server/config"
	"github.com/BlenDMinh/dutgrad-server/database/entity"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var EntityList = []interface{}{
	&entity.User{},
	&entity.UserAuthCredential{},
}

var db *gorm.DB

func connect(driver string, dsn string) *gorm.DB {
	var _db *gorm.DB
	var err error

	_, err = url.ParseRequestURI(dsn)
	if err == nil {
		switch driver {
		case "sqlite":
			_db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		case "mysql":
			_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		case "postgres":
			_db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		default:
			panic(fmt.Sprintf("unsupported database driver: %s", driver))
		}
	} else {
		switch driver {
		case "sqlite":
			_db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		case "mysql":
			_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		case "postgres":
			_db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		default:
			panic(fmt.Sprintf("unsupported database driver: %s", driver))
		}
	}

	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	return _db
}

func Init() {
	config := config.GetEnv()

	if len(config.MasterDBs) == 0 {
		panic("no database configuration found")
	}

	defaultConfig := config.MasterDBs[0]
	driver := defaultConfig.Driver
	dsn := defaultConfig.DSN

	db = connect(driver, dsn)

	if len(config.MasterDBs) > 1 {
		replicas := []gorm.Dialector{}
		for i := 1; i < len(config.MasterDBs); i++ {
			replicaConfig := config.MasterDBs[i]
			driver := replicaConfig.Driver
			dsn := replicaConfig.DSN

			replica := connect(driver, dsn)
			replicas = append(replicas, replica.Dialector)
		}

		db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
	}

	for _, entity := range EntityList {
		db.AutoMigrate(entity)
	}
}

func GetDB() *gorm.DB {
	return db
}
