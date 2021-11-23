package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"home-cloud/utils"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:char(36);primaryKey;"`
	Username string    `gorm:"type:varchar(50);unique;not null"`
	Nickname string    `gorm:"type:varchar(50);not null"`
	Email    string    `gorm:"type:varchar(50)"`
	// 0 for male, 1 for female, 2 for other
	Gender      int    `gorm:"type:tinyint;default:0"`
	Bio         string `gorm:"type:text;default:null"`
	Password    string `gorm:"size:128;not null"`
	AccountSalt string `gorm:"size:64;not null"`
	MacSalt     string `gorm:"size:64;not null"`
	// 0 for user, 1 for admin
	Status int `gorm:"type:tinyint;default:0;comment:'user status"`
	// default 1G quota
	Storage uint64 `gorm:"default:1073741824;comment:'user Storage"`
	// Used storage
	UsedStorage uint64 `gorm:"default:0;comment:'user Storage"`
	// 0 for disable encryption, 1 for AES-256-GCM, 2 for ChaCha20-Poly1305, 3 for XChaCha20-Poly1305
	Encryption    int    `gorm:"type:tinyint;default:0"`
	EncryptionKey string `gorm:"size:64;default:null"`
	Migration     int    `gorm:"type:tinyint;default:0"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.New()
	return nil
}

func (user *User) GetRootFolder() (*File, error) {
	var file File
	err := DB.Where(&File{OwnerId: user.ID}).Where("parent_id = ?", uuid.Nil).First(&file).Error
	return &file, err
}

func GetUserMacSalt(username string) (string, error) {
	var user User
	err := DB.Select("account_salt").Where(&User{Username: username}).First(&user).Error
	return user.AccountSalt, err
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

func GetUserByID(uid uuid.UUID) (*User, error) {
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
	rootFolder.ID = uuid.New()
	rootFolder.ParentId = uuid.Nil
	rootFolder.OwnerId = user.ID
	rootFolder.CreatorId = user.ID
	rootFolder.IsDir = 1
	rootFolder.Name = "Home"
	err = rootFolder.CreateFile()
	return err
}

func (user *User) ChangePassword(newPass string, newAccountSalt string, newMacSalt string) {
	utils.GetLogger().Warn("Change Password for " + user.Username)
	user.Password = newPass
	user.AccountSalt = newAccountSalt
	user.MacSalt = newMacSalt
	DB.Save(&user)
}

func (user *User) UpdateProfile(email string, nickName string, gender int, bio string) {
	utils.GetLogger().Info("Update profile for user " + user.Username)
	user.Email = email
	user.Nickname = nickName
	user.Gender = gender
	user.Bio = bio
	DB.Save(&user)
}

func (user *User) SearchFiles(keyword string) ([]*File, error) {
	var files []*File
	var err error
	err = DB.Model(&File{}).Where(&File{OwnerId: user.ID}).
		Where("name like ?", "%"+keyword+"%").
		Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (user *User) FindFavorites() ([]*File, error) {
	var files []*File
	var err error
	err = DB.Model(&File{}).Where(&File{OwnerId: user.ID, Favorite: 1}).
		Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (user *User) UpdateUsedStorage(newSize uint64) {
	user.UsedStorage = newSize
	DB.Save(user)
}

func (user *User) SetStorageQuota(newSize uint64) {
	user.Storage = newSize
	DB.Save(user)
}

func (user *User) SetAsAdmin() {
	user.Status = 1
	DB.Save(user)
}

func (user *User) SetAsNormalUser() {
	user.Status = 0
	DB.Save(user)
}

func (user *User) DeleteUser() {
	DB.Unscoped().Delete(user)
}

// GetUserList will not include the user who called this method
func (user *User) GetUserList() (users []*User, err error) {
	err = DB.Not(&User{ID: user.ID}).Order("status desc").Order("username").Find(&users).Error
	return
}

// GetAdminCount count admin
func GetAdminCount() (count int64) {
	DB.Model(&User{}).Where(&User{Status: 1}).Count(&count)
	return
}

func (user *User) SetPassword(newPass string, newAccountSalt string, newMacSalt string) {
	user.Password = newPass
	user.AccountSalt = newAccountSalt
	user.MacSalt = newMacSalt
	DB.Save(user)
}
