package gophbot

import (
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	// This is just a convenience import that imports the mysql driver.
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

// State keeps commonly use settings for the bot in memory for quick access
var State = struct {
	sync.RWMutex
	Guilds map[string]*Guild
}{
	Guilds: make(map[string]*Guild),
}

// Guild is the metadata of this guild in the database, for guild-specific settings.
type Guild struct {
	ID     Snowflake `gorm:"size:20;primary;not null"`
	Prefix string    `gorm:"size:10;not null;default:\"!\""`
}
