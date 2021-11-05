package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"home-cloud/models"
	"home-cloud/utils"
	"mime/multipart"
	"os"
	"path"
)

func UploadFile(upFile *multipart.FileHeader, user *models.User, folder *models.File, c *gin.Context) (err error) {
	if folder.OwnerId != user.ID {
		return ErrInvalidOrPermission
	}
	if folder.IsDir != 1 {
		return ErrRequestPara
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
	if err = c.SaveUploadedFile(upFile, dst); err != nil {
		return ErrSave
	}

	err = file.CreateFile()
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			// Duplicate entry error, try to update file
			file, err = updateFile(upFile, user, folder.ID)
			if err != nil {
				return err
			}
		} else {
			return ErrSave
		}
	}
	return nil
}

//Update files when detected duplicate entry in uploading process
func updateFile(upFile *multipart.FileHeader, user *models.User, folderID uuid.UUID) (*models.File, error) {
	file, err := models.GetFileByName(upFile.Filename, user, folderID)
	if err != nil {
		return nil, ErrFoundFile
	}
	if file.IsDir == 1 {
		return nil, ErrConflict
	}
	file.Size = uint64(upFile.Size)
	err = file.UpdateFile()
	if err != nil {
		err = ErrSave
	}
	return file, err
}

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

func NewFileOrFolder(folder *models.File, user *models.User, newName string, t string) (err error) {
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
	file.Size = 0
	file.ParentId = folder.ID
	if file.IsDir == 0 {
		file.FileType = utils.GetFileTypeByName(file.Name)
	}

	if t == "file" {
		dst := path.Join(utils.GetConfig().UserDataPath, user.ID.String(),
			"data", "files", file.RealPath)
		utils.GetLogger().Infof("Create file to %s", dst)
		var tmpFile *os.File
		tmpFile, err = os.Create(dst)
		if err != nil {
			return ErrSystem
		} else {
			err = tmpFile.Close()
			if err != nil {
				return ErrSystem
			}
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

func DeleteFile(file *models.File, user *models.User) (err error) {
	if file.OwnerId != user.ID {
		err = ErrInvalidOrPermission
		return
	}
	//Will not raise error
	DeleteFileRecursively(file, user)
	return nil
}

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
