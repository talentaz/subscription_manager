package db

import (
	"fmt"
	"log"
	"subscriptionManager/models"
	"subscriptionManager/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
        "gorm.io/gorm/schema"
        "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	user := config.DBUsername
	password := config.DBPassword
	host := config.Host
	port := config.Port
	dbname := config.DBName
        dbSchema := config.DBSchema

	dbURL := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
        db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
                     NamingStrategy: schema.NamingStrategy{
                       TablePrefix: dbSchema + ".",
                       NoLowerCase: false,
                     },
                     Logger: logger.Default.LogMode(logger.Info),
                   })

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.AvailablePlans{})
	DB = db
}
