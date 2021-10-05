package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"home-cloud/utils"
	"os"
	"path"
	"strconv"
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
	err := os.MkdirAll(utils.GetConfig().UserDataPath, 0666)
	if err != nil {
		utils.GetLogger().Panic("Create user data path error")
	}
	err = DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &File{})
	if err != nil {
		utils.GetLogger().Panic("Migrate Tables Error")
	}
	if !CheckAdminExist() {
		utils.GetLogger().Info("No admin user, create one......")
		if err = InitAdminUser(); err != nil {
			utils.GetLogger().Panic("Create admin user error!")
		}
		user, err := GetUserByUsername("admin")
		if err != nil {
			utils.GetLogger().Panic("Create admin user error")
		}
		userID := strconv.FormatUint(user.ID, 10)
		if err = os.MkdirAll(path.Join(utils.GetConfig().UserDataPath, userID), 0666); err != nil {
			utils.GetLogger().Panic("Create user folder error")
		}
		if err = os.MkdirAll(path.Join(utils.GetConfig().UserDataPath, userID, "data"), 0666); err != nil {
			utils.GetLogger().Panic("Create user data folder error")
		}
		if err = os.MkdirAll(path.Join(utils.GetConfig().UserDataPath, userID, "data", "files"), 0666); err != nil {
			utils.GetLogger().Panic("Create user file folder error")
		}
	}
}
