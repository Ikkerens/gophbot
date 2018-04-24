package gophbot

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wantedly/gorm-zap"
	"go.uber.org/zap"
)

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", os.Getenv("SQL_DSN")+"?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		return nil, err
	}
	db.LogMode(false)
	db.SetLogger(gormzap.New(Log))

	if err = db.AutoMigrate(new(Guild)).Error; err != nil {
		Log.Fatal("Could not migrate database", zap.Error(err))
	}

	return db, nil
}

type Guild struct {
	ID     Snowflake `gorm:"size:20;primary;not null"`
	Prefix string    `gorm:"size:10;not null;default:\"!\""`
}
