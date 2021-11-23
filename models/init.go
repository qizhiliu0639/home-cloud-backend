package models

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"home-cloud/utils"
	"os"
	"path"
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

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	DB = db
	Migration()
}

func Migration() {
	err := os.MkdirAll(utils.GetConfig().UserDataPath, 0644)
	if err != nil {
		panic("Create user data path error: " + err.Error())
	}
	err = DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &File{})
	if err != nil {
		panic("Migrate tables error: " + err.Error())
	}
	if !CheckAdminExist() {
		utils.GetLogger().Info("No admin user, create one......")
		if err = InitAdminUser(); err != nil {
			panic("Create admin user error: " + err.Error())
		}
		var user *User
		user, err = GetUserByUsername("admin")
		if err != nil {
			panic("Create admin user error: " + err.Error())
		}
		userID := user.ID.String()
		if err = os.MkdirAll(path.Join(utils.GetConfig().UserDataPath, userID), 0644); err != nil {
			panic("Create user folder error: " + err.Error())
		}
		if err = os.MkdirAll(path.Join(utils.GetConfig().UserDataPath, userID, "data"), 0644); err != nil {
			panic("Create user data folder error: " + err.Error())
		}
		if err = os.MkdirAll(path.Join(utils.GetConfig().UserDataPath, userID, "data", "files"), 0644); err != nil {
			panic("Create user file folder error: " + err.Error())
		}
	}
}

func CheckAdminExist() bool {
	var user User
	err := DB.Where(&User{Status: 1}).First(&user).Error
	if err != nil {
		return false
	} else {
		return true
	}
}

func InitAdminUser() error {
	adminUser := NewUser()
	adminUser.ID = uuid.New()
	adminUser.Username = "admin"
	adminUser.Nickname = "admin"
	newAccountSalt, newMacSalt, newSavePassword, err := utils.GeneratePasswordInfoFromPassword("admin")
	if err != nil {
		return err
	}
	adminUser.AccountSalt = newAccountSalt
	adminUser.MacSalt = newMacSalt
	adminUser.Password = newSavePassword
	adminUser.Status = 1
	err = adminUser.RegisterUser()
	if err != nil {
		return err
	} else {
		fmt.Println("Creating admin user completes.")
		fmt.Println("Username: admin")
		fmt.Println("Password: admin")
		fmt.Println("This notification will not show again. Please save the username and password.")
		fmt.Println("Press enter to continue...")
		_, _ = fmt.Scanln()
	}
	return nil
}
