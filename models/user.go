package models

import (
	"gorm.io/gorm"
	"home-cloud/utils"
)

type User struct {
	gorm.Model
	ID          uint64 `gorm:"primaryKey;AUTO_INCREMENT"`
	Username    string `gorm:"type:varchar(50);unique;not null"`
	Nickname    string `gorm:"type:varchar(50);not null"`
	Email       string `gorm:"type:varchar(50)"`
	Password    string `gorm:"size:512;not null"`
	AccountSalt string `gorm:"size:256;not null"`
	MacSalt     string `gorm:"size:256;not null"`
	// 0 for user, 1 for admin, 2 for resetting password, 3 for resetting two-factor auth
	// 4 for resetting both, 5 for disabled user
	Status  int `gorm:"default:0;comment:'user status"`
	Storage uint64
}

func (user *User) GetRootFolder() (*File, error) {
	var file File
	err := DB.Where(&File{OwnerId: user.ID}).
		Where("parent_id is null").First(&file).Error
	return &file, err
}

func GetUserMacSalt(username string) (string, error) {
	var user User
	err := DB.Select("mac_salt").Where(&User{Username: username}).First(&user).Error
	return user.MacSalt, err
}

func GetUserPassword(username string) (User, error) {
	var user User
	err := DB.Select("password", "account_salt").
		Where(&User{Username: username}).First(&user).Error
	return user, err
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where(&User{Username: username}).First(&user).Error
	return &user, err
}

func GetUserByID(uid uint64) (*User, error) {
	var user User
	err := DB.Where(&User{ID: uid}).First(&user).Error
	return &user, err
}

func NewUser() *User {
	return &User{}
}

func (user *User) RegisterUser() error {
	utils.GetLogger().Info(user.Status)
	err := DB.Create(user).Error
	if err != nil {
		return err
	}
	rootFolder := NewFile()
	rootFolder.OwnerId = user.ID
	rootFolder.CreatorId = user.ID
	rootFolder.IsDir = 1
	rootFolder.Name = user.Username
	err = rootFolder.CreateFile()
	return err
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
	adminUser.Username = "admin"
	adminUser.Nickname = "admin"
	adminUser.AccountSalt = utils.GenerateSalt(256)
	adminUser.Password = utils.GetHashWithSalt("admin", adminUser.AccountSalt)
	adminUser.MacSalt = utils.GenerateSalt(256)
	adminUser.Status = 1
	err := adminUser.RegisterUser()
	return err
}
