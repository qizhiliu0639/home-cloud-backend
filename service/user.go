package service

import (
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
		macSalt := utils.GenerateSalt(256)
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
	newMacSalt := utils.GenerateSalt(256)
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

func AddManageUserAuth(AdminUser *models.User, user *models.User) (err error) {
	if AdminUser.Status != 1 {
		err = ErrInvalidAuth
		return err
	}

	err = user.AddMangerUserAuth()
	if err != nil {
		return err
	}
	return nil
}

func CancelManageUserAuth(AdminUser *models.User, user *models.User) (err error) {
	if AdminUser.Status != 1 {
		err = ErrInvalidAuth
		return err
	}

	err = user.CancelMangerUserAuth()
	if err != nil {
		return err
	}
	return nil
}

func AdjustUserStorage(AdminUser *models.User, user *models.User) (err error) {
	if AdminUser.Status != 1 && AdminUser.Status != 6 {
		err = ErrInvalidAuth
		return err
	}

	err = user.UpdateStorage(user.Storage + 10240)
	if err != nil {
		err = ErrStorage
	}
	return err
}
