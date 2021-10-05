package service

import (
	"errors"
	"home-cloud/models"
	"home-cloud/utils"
	"os"
	"path"
	"strconv"
)

func LoginGetSalt(username string) string {
	salt, err := models.GetUserMacSalt(username)
	if err != nil {
		salt = utils.GenerateSalt(256)
	}
	return salt
}

func LoginValidate(username string, password string) bool {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return false
	}
	if utils.GetHashWithSalt(password, user.AccountSalt) != user.Password {
		return false
	}
	return true
}

func RegisterUser(username string, password string, macSalt string) error {
	if _, err := models.GetUserByUsername(username); err != nil {
		user := models.NewUser()
		user.Username = username
		accountSalt := utils.GenerateSalt(256)
		user.AccountSalt = accountSalt
		user.Password = utils.GetHashWithSalt(password, accountSalt)
		user.MacSalt = macSalt
		user.Nickname = username
		err := user.RegisterUser()
		if err != nil {
			return err
		}
		//Create user folder
		user, err = models.GetUserByUsername(username)
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
		return nil
	} else {
		return errors.New("username unavailable")
	}
}
