package service

import (
	"home-cloud/models"
	"home-cloud/utils"
	"io/ioutil"
	"os"
	"path"
)

func MigrateAlgorithm(user *models.User, oldAlgorithm int, newAlgorithm int, fileEncryptionKey []byte) {
	utils.GetLogger().Info("Migrating encryption algorithm for user " + user.Username)
	user.SetEncryption(newAlgorithm)
	userFilePath := path.Join(utils.GetConfig().UserDataPath, user.ID.String(),
		"data", "files")
	files, err := os.ReadDir(userFilePath)
	if err != nil {
		utils.GetLogger().Fatal("Migration encrypted file for user " + user.Username + " error: " + err.Error())
		user.SetMigration(2)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := path.Join(userFilePath, file.Name())
		var fileContent []byte
		fileContent, err = ioutil.ReadFile(filePath)
		if err != nil {
			utils.GetLogger().Error("Read file " + filePath + " for user " + user.Username + " error: " + err.Error())
			continue
		}
		var originContent []byte
		var errorDecrypt error
		if oldAlgorithm == 1 {
			originContent, errorDecrypt = utils.DecryptFileAES(fileEncryptionKey, fileContent)
		} else if oldAlgorithm == 2 {
			originContent, errorDecrypt = utils.DecryptFileChaCha(fileEncryptionKey, fileContent)
		} else if oldAlgorithm == 3 {
			originContent, errorDecrypt = utils.DecryptFileXChaCha(fileEncryptionKey, fileContent)
		} else {
			originContent = fileContent
		}
		if errorDecrypt != nil {
			utils.GetLogger().Error("Decrypt file " + filePath + " for user " + user.Username + " error: " + errorDecrypt.Error())
			continue
		}
		var newContent []byte
		var errorEncrypt error
		if newAlgorithm == 1 {
			newContent, errorEncrypt = utils.EncryptFileAES(fileEncryptionKey, originContent)
		} else if newAlgorithm == 2 {
			newContent, errorEncrypt = utils.EncryptFileChaCha(fileEncryptionKey, originContent)
		} else if newAlgorithm == 3 {
			newContent, errorEncrypt = utils.EncryptFileXChaCha(fileEncryptionKey, originContent)
		} else {
			newContent = originContent
		}
		if errorEncrypt != nil {
			utils.GetLogger().Error("Encrypt file " + filePath + " for user " + user.Username + " error: " + errorEncrypt.Error())
			continue
		}

		err = ioutil.WriteFile(filePath, newContent, 0644)
		if err != nil {
			utils.GetLogger().Error("Write file " + filePath + " for user " + user.Username + " error: " + err.Error())
			continue
		}
	}
	utils.GetLogger().Info("Migrating encryption algorithm for user " + user.Username + " completes")
	user.SetMigration(0)
}
