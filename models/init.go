package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"home-cloud/utils"
)

var DB *gorm.DB

func InitDatabase() {
	config := utils.GetConfig()
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	utils.GetLogger().Infof("Try to Connect to %s", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		utils.GetLogger().Panic("Failed to connect to the database")
	}
	DB = db
	Migration()
}

func Migration() {
	err := DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &File{})
	if err != nil {
		utils.GetLogger().Panic("Migrate Tables Error")
	}
	if !CheckAdminExist() {
		utils.GetLogger().Info("No admin user, create one......")
		if err = InitAdminUser(); err != nil {
			utils.GetLogger().Panic("Create admin user error!")
		}
	}
}
