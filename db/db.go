package db

import (
	"github.com/adigunhammedolalekan/sms-forwarder/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Connect(uri string) (*gorm.DB, error) {
	conn, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	runMigration(conn)
	return conn, nil
}

func runMigration(conn *gorm.DB) {
	conn.Debug().AutoMigrate(&types.User{})
}
