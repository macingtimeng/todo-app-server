package db

import (
	"fmt"
	"log"
	"time"
	"todo-app/entity"
	"todo-app/infra/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DbConn() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.AppConfig().DbHost,
		config.AppConfig().DbUser,
		config.AppConfig().DbPassword,
		config.AppConfig().DbName,
		config.AppConfig().DbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Panic("error while connect to db: ", err.Error())
	}

	// connection pool
	d, err := db.DB()

	if err != nil {
		log.Panic(err.Error())
	}

	d.SetConnMaxIdleTime(10 * time.Minute)
	d.SetConnMaxLifetime(1 * time.Hour)
	d.SetMaxIdleConns(10)
	d.SetMaxOpenConns(100)

	err = db.AutoMigrate(&entity.User{}, &entity.Todo{})

	if err != nil {
		log.Panic("error while migration: ", err.Error())
	}

	return db
}
