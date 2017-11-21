package model

import (
	"fmt"

	"github.com/dmfrank/egaas/config"
	"github.com/jinzhu/gorm"
)

var (
	// DBConn connection handler
	DBConn *gorm.DB
)

// GormInit initialising db connection
func GormInit() error {
	conf := &config.DBConfig{}
	err := conf.Read()
	if err != nil {
		fmt.Printf("Configuration file reading issue: %s", err.Error())
		return err
	}

	DBConn, err = gorm.Open("postgres",
		fmt.Sprintf(
			"host=localhost user=%s dbname=%s sslmode=disable password=%s",
			conf.DBUser,
			conf.DBName,
			conf.DBPass))
	if err != nil {
		fmt.Printf("DB connection error: %s", err.Error())
		return err
	}
	return nil
}

// GormClose closing db connection
func GormClose() error {
	if DBConn != nil {
		return DBConn.Close()
	}
	return nil
}
