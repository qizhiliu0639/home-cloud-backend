package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"home-cloud/models"
	"home-cloud/utils"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

// UploadFile upload file to the folder
func UploadFile(upFile *multipart.FileHeader, user *models.User, folder *models.File, c *gin.Context) (err error) {
	if folder.OwnerId != user.ID {
		return ErrInvalidOrPermission
	}
	if folder.IsDir != 1 {
		return ErrRequestPara
	}
	if user.UsedStorage+uint64(upFile.Size) > user.Storage {
		return ErrStorage
	}
	file := models.NewFile()
	file.ID = uuid.New()
	file.RealPath = file.ID.String()
	file.IsDir = 0
	file.Name = upFile.Filename
	file.OwnerId = user.ID
	file.CreatorId = user.ID
	file.Size = uint64(upFile.Size)
	file.ParentId = folder.ID
	file.FileType = utils.GetFileTypeByName(file.Name)

	dst := path.Join(utils.GetConfig().UserDataPath, user.ID.String(),
		"data", "files", file.RealPath)
	utils.GetLogger().Infof("Save file to %s", dst)

	if user.Encryption > 3 || user.Encryption < 0 {
		return ErrSystem
	}
	if err = saveUploadFileEncryption(upFile, dst, user, c); err != nil {
		return err
	}
	err = file.CreateFile()
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			// Duplicate entry error, try to update file
			file, err = updateFile(upFile, user, folder.ID, dst)
			if err != nil {
				return err
			}
		} else {
			return ErrSave
		}
	}
	user.UpdateUsedStorage(user.UsedStorage + file.Size)
	return nil
}

//Update files when detected duplicate entry in uploading process
func updateFile(upFile *multipart.FileHeader, user *models.User, folderID uuid.UUID, newFilePath string) (*models.File, error) {
	file, err := models.GetFileByName(upFile.Filename, user, folderID)
	if err != nil {
		return nil, ErrFoundFile
	}
	if file.IsDir == 1 {
		return nil, ErrConflict
	}
	oldSize := file.Size
	oldFilePath := path.Join(utils.GetConfig().UserDataPath, user.ID.String(),
		"data", "files", file.RealPath)
	// Will replace the old file
	err = os.Rename(newFilePath, oldFilePath)
	if err != nil {
		return nil, ErrSystem
	}
	file.Size = uint64(upFile.Size)
	err = file.UpdateFile()
	if err != nil {
		err = ErrSave
	}
	user.UpdateUsedStorage(user.UsedStorage - oldSize)
	return file, err
}

// saveUploadFileEncryption will save the upload file to the local file system
// If user setting encryption is enabled, it will encrypt the file before writing to the system
func saveUploadFileEncryption(upFile *multipart.FileHeader, dst string, user *models.User, c *gin.Context) error {
	encryptedKey := c.Value("encryptionKey").([]byte)
	fileEncryptionKey, err := utils.DecryptEncryptionKey(encryptedKey, user.EncryptionKey)
	if err != nil {
		return ErrRequestPara
	}
	var file multipart.File
	file, err = upFile.Open()
	if err != nil {
		return ErrRequestPara
	}
	var fileContent []byte
	fileContent, err = ioutil.ReadAll(file)
	errClose := file.Close()
	if err != nil || errClose != nil {
		return ErrRequestPara
	}
	var encryptedContent []byte
	var errEncrypt error
	if user.Encryption == 1 {
		encryptedContent, errEncrypt = utils.EncryptFileAES(fileEncryptionKey, fileContent)
	} else if user.Encryption == 2 {
		encryptedContent, errEncrypt = utils.EncryptFileChaCha(fileEncryptionKey, fileContent)
	} else if user.Encryption == 3 {
		encryptedContent, errEncrypt = utils.EncryptFileXChaCha(fileEncryptionKey, fileContent)
	} else {
		encryptedContent = fileContent
	}
	if errEncrypt != nil {
		return ErrSystem
	}
	err = ioutil.WriteFile(dst, encryptedContent, 0644)
	if err != nil {
		return ErrSave
	}
	return nil
}

// GetFolder return children in the folder
func GetFolder(folder *models.File, user *models.User) (files []*models.File, err error) {
	if folder.IsDir != 1 {
		return nil, ErrRequestPara
	}
	if folder.OwnerId != user.ID {
		return nil, ErrInvalidOrPermission
	}
	err = folder.TraceRoot()
	if err != nil {
		return nil, ErrSystem
	}
	files, err = folder.GetChildInFolder()
	if err != nil {
		return nil, ErrSystem
	}
	return files, err
}

// NewFileOrFolder create a file or a folder in the current folder
func NewFileOrFolder(folder *models.File, user *models.User, newName string, t string, c *gin.Context) (err error) {
	if folder.OwnerId != user.ID {
		return ErrInvalidOrPermission
	}
	if folder.IsDir != 1 {
		return ErrRequestPara
	}
	file := models.NewFile()
	if t == "file" {
		file.IsDir = 0
	} else if t == "folder" {
		file.IsDir = 1
	} else {
		return ErrRequestPara
	}
	file.ID = uuid.New()
	file.RealPath = file.ID.String()
	file.Name = newName
	file.OwnerId = user.ID
	file.CreatorId = user.ID
	// new file or folder size will be always 0, no need to update UsedStorage
	file.Size = 0
	file.ParentId = folder.ID
	if file.IsDir == 0 {
		file.FileType = utils.GetFileTypeByName(file.Name)
	}

	if t == "file" {
		dst := path.Join(utils.GetConfig().UserDataPath, user.ID.String(),
			"data", "files", file.RealPath)
		utils.GetLogger().Infof("Create file to %s", dst)
		// If the user encryption setting is enabled, it will also encrypt the empty file
		encryptedKey := c.Value("encryptionKey").([]byte)
		var fileEncryptionKey []byte
		fileEncryptionKey, err = utils.DecryptEncryptionKey(encryptedKey, user.EncryptionKey)
		if err != nil {
			return ErrRequestPara
		}
		fileContent := make([]byte, 0)
		var encryptedContent []byte
		var errEncrypt error
		if user.Encryption == 1 {
			encryptedContent, errEncrypt = utils.EncryptFileAES(fileEncryptionKey, fileContent)
		} else if user.Encryption == 2 {
			encryptedContent, errEncrypt = utils.EncryptFileChaCha(fileEncryptionKey, fileContent)
		} else if user.Encryption == 3 {
			encryptedContent, errEncrypt = utils.EncryptFileXChaCha(fileEncryptionKey, fileContent)
		} else {
			encryptedContent = fileContent
		}
		if errEncrypt != nil {
			return ErrSystem
		}
		err = ioutil.WriteFile(dst, encryptedContent, 0644)
		if err != nil {
			return ErrSystem
		}
	}

	err = file.CreateFile()
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrDuplicate
		} else {
			return ErrSave
		}
	}

	return nil
}

// GetFile return path pointed to requested file in the user data folder
func GetFile(file *models.File, user *models.User) (dst, filename string, err error) {
	if file.OwnerId != user.ID {
		err = ErrInvalidOrPermission
		return
	}
	if file.IsDir != 0 {
		err = ErrRequestPara
		return
	}
	filename = file.Name
	dst = path.Join(utils.GetConfig().UserDataPath, user.ID.String(),
		"data", "files", file.RealPath)
	return
}

// GetFileEncrypted will decrypt the file and return the original file content
func GetFileEncrypted(dst string, user *models.User, c *gin.Context) ([]byte, error) {
	encryptedKey := c.Value("encryptionKey").([]byte)
	fileEncryptionKey, err := utils.DecryptEncryptionKey(encryptedKey, user.EncryptionKey)
	if err != nil {
		return nil, ErrRequestPara
	}
	var encryptedFile []byte
	encryptedFile, err = ioutil.ReadFile(dst)
	if err != nil {
		return nil, ErrSystem
	}
	var decryptedContent []byte
	if user.Encryption == 1 {
		decryptedContent, err = utils.DecryptFileAES(fileEncryptionKey, encryptedFile)
	} else if user.Encryption == 2 {
		decryptedContent, err = utils.DecryptFileChaCha(fileEncryptionKey, encryptedFile)
	} else if user.Encryption == 3 {
		decryptedContent, err = utils.DecryptFileXChaCha(fileEncryptionKey, encryptedFile)
	}
	if err != nil {
		return nil, ErrSystem
	}
	return decryptedContent, nil

}

// GetFileOrFolderInfoByPath return file or folder
func GetFileOrFolderInfoByPath(paths []string, user *models.User) (*models.File, error) {
	rootFolder, err := user.GetRootFolder()
	if err != nil {
		return nil, ErrSystem
	}
	//root folder
	if len(paths) == 0 {
		err = rootFolder.TraceRoot()
		if err != nil {
			return nil, ErrSystem
		}
		return rootFolder, nil
	}
	var file *models.File
	for _, p := range paths {
		file, err = rootFolder.GetChildInFolderByName(p)
		if err != nil || file.OwnerId != user.ID {
			return nil, ErrInvalidOrPermission
		}
		rootFolder = file
	}
	err = file.TraceRoot()
	if err != nil {
		return nil, ErrSystem
	}
	return file, nil
}

// GetFileOrFolderInfoByID return file or folder
func GetFileOrFolderInfoByID(vFileID uuid.UUID, user *models.User) (*models.File, error) {
	file, err := models.GetFileByID(vFileID)
	if err != nil {
		return nil, ErrInvalidOrPermission
	}
	if file.OwnerId != user.ID {
		return nil, ErrInvalidOrPermission
	}
	err = file.TraceRoot()
	if err != nil {
		return nil, ErrSystem
	}
	return file, nil
}

// DeleteFile delete a folder or file
func DeleteFile(file *models.File, user *models.User) (err error) {
	if file.OwnerId != user.ID {
		err = ErrInvalidOrPermission
		return
	}
	//Will not raise error
	DeleteFileRecursively(file, user)
	return nil
}

// DeleteFileRecursively help to delete a file or folder recursively
func DeleteFileRecursively(file *models.File, user *models.User) {
	deleteQueue := []*models.File{file}
	//Max 65536 level
	level := 0
	count := 1
	for len(deleteQueue) > 0 && level < 65536 {
		count--
		root := deleteQueue[0]
		if root.IsDir == 1 {
			child, err := root.GetChildInFolder()
			if err == nil {
				deleteQueue = append(deleteQueue, child...)
			}
		} else {
			// Reduce used storage
			user.UpdateUsedStorage(user.UsedStorage - root.Size)
		}
		root.DeleteFile()
		dst := path.Join(utils.GetConfig().UserDataPath, user.ID.String(),
			"data", "files", root.RealPath)
		err := os.Remove(dst)
		utils.GetLogger().Info("Delete file in " + dst)
		//Will skip deleting the file if error
		if err != nil {
			utils.GetLogger().Error("Error deleting " + dst)
		}
		deleteQueue = deleteQueue[1:]

		if count == 0 {
			level++
			count = len(deleteQueue)
		}
	}
}

// ChangeFavoriteStatus change the favorite setting in the system
func ChangeFavoriteStatus(file *models.File, user *models.User) (err error) {
	if file.OwnerId != user.ID {
		err = ErrInvalidOrPermission
		return
	}
	if file.Favorite == 0 {
		err = file.AddFavorite()
	} else {
		err = file.CancelFavorite()
	}
	if err != nil {
		err = ErrFavorite
	}
	return err
}

// GetFavorites return files and folders that are set favorite
func GetFavorites(user *models.User) (files []*models.File, err error) {
	files, err = user.FindFavorites()
	if err != nil {
		err = ErrSystem
	}
	for _, v := range files {
		err = v.TraceRoot()
		if err != nil {
			err = ErrSystem
			return
		}
	}
	return
}

// SearchFiles return files based on keyword
func SearchFiles(user *models.User, keyword string) (files []*models.File, err error) {
	files, err = user.SearchFiles(keyword)
	if err != nil {
		err = ErrSystem
		return
	}
	for _, v := range files {
		err = v.TraceRoot()
		if err != nil {
			err = ErrSystem
			return
		}
	}
	return
}
