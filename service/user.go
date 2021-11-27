package service

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"home-cloud/models"
	"home-cloud/utils"
	"os"
	"path"
	"strings"
)

// LoginGetSalt return the salt of the user or a fake salt if the user not exists
func LoginGetSalt(username string) string {
	accountSalt, err := models.GetUserMacSalt(username)
	if err != nil {
		accountSalt = utils.GenerateFakeSalt(username)
	}
	return accountSalt
}

// LoginValidate validate if the username and password match
func LoginValidate(username string, password string) (bool, *models.User) {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return false, nil
	}
	if utils.GetHashWithSalt(password, user.MacSalt) != user.Password {
		return false, nil
	}
	return true, user
}

// RegisterUser register a user in the system
func RegisterUser(username string, password string, accountSalt string, encryptionKey string) error {
	if _, err := models.GetUserByUsername(username); err != nil {
		user := models.NewUser()
		user.Username = username
		user.AccountSalt = accountSalt
		macSalt := utils.GenerateSaltOrKey()
		user.MacSalt = macSalt
		user.Password = utils.GetHashWithSalt(password, macSalt)
		user.Nickname = username
		var encryptionKeyByte []byte
		encryptionKeyByte, err = hex.DecodeString(encryptionKey)
		if err != nil {
			return ErrRequestPara
		}
		var encryptKey []byte
		encryptKey, err = hex.DecodeString(utils.GenerateSaltOrKey())
		if err != nil {
			return err
		}
		var newEncryptionKey string
		newEncryptionKey, err = utils.EncryptEncryptionKey(encryptionKeyByte, encryptKey)
		if err != nil {
			utils.GetLogger().Panic("Create user error")
			return err
		}
		user.EncryptionKey = newEncryptionKey
		err = user.RegisterUser()
		if err != nil {
			utils.GetLogger().Panic("Create user error")
			return err
		}
		//Create user folder
		user, err = models.GetUserByUsername(username)
		if err != nil {
			utils.GetLogger().Panic("Create user error")
			return err
		}
		userID := user.ID.String()
		if err = os.MkdirAll(path.Join(utils.GetConfig().UserDataPath, userID), 0666); err != nil {
			utils.GetLogger().Panic("Create user folder error")
			return err
		}
		if err = os.MkdirAll(path.Join(utils.GetConfig().UserDataPath, userID, "data"), 0666); err != nil {
			utils.GetLogger().Panic("Create user data folder error")
			return err
		}
		if err = os.MkdirAll(path.Join(utils.GetConfig().UserDataPath, userID, "data", "files"), 0666); err != nil {
			utils.GetLogger().Panic("Create user file folder error")
			return err
		}
		return nil
	} else {
		return ErrUsernameInvalid
	}
}

// ChangePassword change user password
func ChangePassword(user *models.User, newAccountSalt string, newPassword string, oldEncryption string, newEncryption string) error {
	newMacSalt := utils.GenerateSaltOrKey()
	newPass := utils.GetHashWithSalt(newPassword, newMacSalt)
	var newEncryptionKeyByte []byte
	var err error
	newEncryptionKeyByte, err = hex.DecodeString(newEncryption)
	if err != nil {
		return ErrRequestPara
	}
	var oldEncryptionKeyByte []byte
	oldEncryptionKeyByte, err = hex.DecodeString(oldEncryption)
	if err != nil {
		return ErrRequestPara
	}
	var fileEncryptionKey []byte
	fileEncryptionKey, err = utils.DecryptEncryptionKey(oldEncryptionKeyByte, user.EncryptionKey)
	if err != nil {
		return ErrRequestPara
	}
	var newFileEncryptionKey string
	newFileEncryptionKey, err = utils.EncryptEncryptionKey(newEncryptionKeyByte, fileEncryptionKey)
	if err != nil {
		return ErrRequestPara
	}
	user.ChangePassword(newPass, newAccountSalt, newMacSalt, newFileEncryptionKey)
	return nil
}

// UpdateProfile update user profile settings in the database
func UpdateProfile(user *models.User, email string, nickName string, gender int, bio string) (err error) {
	if gender < 0 || gender > 2 || !strings.Contains(email, "@") {
		return ErrRequestPara
	}
	user.UpdateProfile(email, nickName, gender, bio)
	return nil
}

// GetUserList return current users in the system
func GetUserList(user *models.User) ([]*models.User, error) {
	users, err := user.GetUserList()
	if err != nil {
		return nil, ErrSystem
	}
	// The user will be always be the first element in the result array
	return append([]*models.User{user}, users...), nil
}

// DeleteUser delete the user
func DeleteUser(user *models.User, deleteUserName string) error {
	// Cannot self delete
	if deleteUserName == user.Username {
		return ErrRequestPara
	}
	deleteUser, err := models.GetUserByUsername(deleteUserName)
	if err != nil {
		return ErrRequestPara
	}
	dst := path.Join(utils.GetConfig().UserDataPath, deleteUser.ID.String())
	err = os.RemoveAll(dst)
	if err != nil {
		return ErrSystem
	}
	deleteUser.DeleteUser()
	return nil
}

// ToggleAdmin change user permission
func ToggleAdmin(user *models.User, toggleUsername string) error {
	// Cannot self modify
	if user.Username == toggleUsername {
		return ErrRequestPara
	}
	toggleUser, err := models.GetUserByUsername(toggleUsername)
	if err != nil {
		return ErrRequestPara
	}
	if toggleUser.Status == 0 {
		toggleUser.SetAsAdmin()
	} else {
		if models.GetAdminCount() <= 1 {
			return ErrRequestPara
		} else {
			toggleUser.SetAsNormalUser()
		}
	}
	return nil
}

// SetUserQuota set user storage quota
func SetUserQuota(modifiedUsername string, newSize uint64) error {
	modifiedUser, err := models.GetUserByUsername(modifiedUsername)
	if err != nil {
		return ErrRequestPara
	}
	modifiedUser.SetStorageQuota(newSize)
	return nil
}

// ResetUserPassword reset the user password
func ResetUserPassword(resetUsername string) (string, error) {
	resetUser, err := models.GetUserByUsername(resetUsername)
	if err != nil {
		return "", ErrRequestPara
	}
	if resetUser.Encryption > 0 {
		return "", ErrResetForbidden
	}
	var newPassword, newAccountSalt, newMacSalt, newSavePassword, newEncryptionKey string
	newPassword, newAccountSalt, newMacSalt, newSavePassword, newEncryptionKey, err = utils.GeneratePasswordInfo()
	if err != nil {
		return "", ErrSystem
	}
	resetUser.SetPassword(newSavePassword, newAccountSalt, newMacSalt, newEncryptionKey)
	return newPassword, nil
}

// GetUserNameByID get username by id, return empty string if error
func GetUserNameByID(uid uuid.UUID) string {
	if user, err := models.GetUserByID(uid); err != nil {
		return ""
	} else {
		return user.Username
	}
}

// ChangeEncryptionAlgorithm will set the Migration status of the user and use goroutine to call
// migration process asynchronously
func ChangeEncryptionAlgorithm(user *models.User, algo int, c *gin.Context) error {
	encryptedKey := c.Value("encryptionKey").([]byte)
	fileEncryptionKey, err := utils.DecryptEncryptionKey(encryptedKey, user.EncryptionKey)
	if err != nil {
		return ErrRequestPara
	}
	if algo < 0 || algo > 3 {
		return ErrRequestPara
	}
	user.SetMigration(1)
	go MigrateAlgorithm(user, user.Encryption, algo, fileEncryptionKey)
	return nil
}
