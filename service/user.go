package service

import (
	"fmt"
	"github.com/google/uuid"
	"home-cloud/models"
	"home-cloud/utils"
	"os"
	"path"
	"strings"
)

func LoginGetSalt(username string) string {
	accountSalt, err := models.GetUserMacSalt(username)
	if err != nil {
		accountSalt = utils.GenerateFakeSalt(username)
	}
	return accountSalt
}

func LoginValidate(username string, password string) bool {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return false
	}
	if utils.GetHashWithSalt(password, user.MacSalt) != user.Password {
		return false
	}
	return true
}

func RegisterUser(username string, password string, accountSalt string) error {
	if _, err := models.GetUserByUsername(username); err != nil {
		user := models.NewUser()
		user.Username = username
		user.AccountSalt = accountSalt
		macSalt := utils.GenerateSalt()
		user.MacSalt = macSalt
		user.Password = utils.GetHashWithSalt(password, macSalt)
		user.Nickname = username
		err = user.RegisterUser()
		if err != nil {
			return err
		}
		//Create user folder
		user, err = models.GetUserByUsername(username)
		if err != nil {
			utils.GetLogger().Panic("Create admin user error")
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

func ChangePassword(user *models.User, newAccountSalt string, newPassword string) {
	newMacSalt := utils.GenerateSalt()
	newPass := utils.GetHashWithSalt(newPassword, newMacSalt)
	user.ChangePassword(newPass, newAccountSalt, newMacSalt)
}

func UpdateProfile(user *models.User, email string, nickName string, gender int, bio string) (err error) {
	if gender < 0 || gender > 2 || !strings.Contains(email, "@") {
		return ErrRequestPara
	}
	user.UpdateProfile(email, nickName, gender, bio)
	return nil
}

func GetUserList(user *models.User) ([]*models.User, error) {
	users, err := user.GetUserList()
	if err != nil {
		return nil, ErrSystem
	}
	// The user will be always be the first element in the result array
	return append([]*models.User{user}, users...), nil
}

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
		fmt.Println(err)
		return ErrSystem
	}
	deleteUser.DeleteUser()
	return nil
}

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

func SetUserQuota(modifiedUsername string, newSize uint64) error {
	modifiedUser, err := models.GetUserByUsername(modifiedUsername)
	if err != nil {
		return ErrRequestPara
	}
	modifiedUser.SetStorageQuota(newSize)
	return nil
}

func ResetUserPassword(resetUsername string) (string, error) {
	resetUser, err := models.GetUserByUsername(resetUsername)
	if err != nil {
		return "", ErrRequestPara
	}
	if resetUser.Encryption > 0 {
		return "", ErrResetForbidden
	}
	var newPassword, newAccountSalt, newMacSalt, newSavePassword string
	newPassword, newAccountSalt, newMacSalt, newSavePassword, err = utils.GeneratePasswordInfo()
	if err != nil {
		return "", ErrSystem
	}
	resetUser.SetPassword(newSavePassword, newAccountSalt, newMacSalt)
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
