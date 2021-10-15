package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"path"
)

type File struct {
	// 表字段
	gorm.Model
	ID uuid.UUID `gorm:"primaryKey"`
	//1 for folder and 0 for file
	IsDir int    `gorm:"default:0;not null"`
	Name  string `gorm:"type:varchar(255);uniqueIndex:idx_only_one"`
	//null for root folder
	ParentId  uuid.UUID `gorm:"default:null;uniqueIndex:idx_only_one"`
	OwnerId   uuid.UUID
	CreatorId uuid.UUID
	Size      uint64 `gorm:"not null"`
	//1 for special hidden file or folder, 0 for others
	Status         int    `gorm:"default:0"`
	Locked         int    `gorm:"default:0"`
	Version        uint64 `gorm:"default:0"`
	FileType       int    `gorm:"default:0"`
	RealPath       string `gorm:"not null"`
	Thumbnail      uint64
	SharedEnabled  int `gorm:"default:0"`
	DownloadStatus int `gorm:"default:0"`
	Encryption     int `gorm:"default:0"`

	// 数据库忽略字段
	Position string `gorm:"-"`
}

func (file *File) TraceRoot() (err error) {
	if len(file.Position) > 0 {
		return nil
	}
	if file.ParentId != uuid.Nil {
		var folder *File
		err = DB.Where(&File{ID: file.ParentId, OwnerId: file.OwnerId}).First(&folder).Error
		if err != nil {
			return err
		}
		err = folder.TraceRoot()
		if err != nil {
			return err
		}
		file.Position = path.Join(folder.Position, file.Name)
	} else {
		file.Position = "/"
	}
	return nil
}
func (file *File) CreateFile() error {
	return DB.Create(file).Error
}

func (file *File) UpdateFile() error {
	return DB.Save(file).Error
}

func GetFileByID(fid uuid.UUID) (*File, error) {
	var file File
	err := DB.Where(&File{ID: fid}).First(&file).Error
	return &file, err
}
func GetFileByName(name string, owner *User, parent uuid.UUID) (*File, error) {
	var file File
	err := DB.Where(&File{Name: name, OwnerId: owner.ID, ParentId: parent}).First(&file).Error
	return &file, err
}

func (file *File) GetChildInFolder() ([]*File, error) {
	if file.IsDir == 0 {
		return nil, errors.New("not a folder")
	}
	err := file.TraceRoot()
	if err != nil {
		return nil, err
	}
	var children []*File
	err = DB.Where(&File{ParentId: file.ID, OwnerId: file.OwnerId}).Find(&children).Error
	if (err != nil) && (!errors.Is(err, gorm.ErrRecordNotFound)) {
		return nil, err
	}
	for _, child := range children {
		child.Position = path.Join(file.Position, child.Name)
	}
	return children, nil
}

func (file *File) GetChildInFolderByName(filename string) (*File, error) {
	if file.IsDir == 0 {
		return nil, errors.New("not a folder")
	}
	err := file.TraceRoot()
	if err != nil {
		return nil, err
	}
	var child *File
	err = DB.Where(&File{ParentId: file.ID, Name: filename, OwnerId: file.OwnerId}).First(&child).Error
	if (err != nil) && (!errors.Is(err, gorm.ErrRecordNotFound)) {
		return nil, err
	}
	child.Position = path.Join(file.Position, child.Name)
	return child, nil
}

func NewFile() *File {
	return &File{}
}

func (file *File) DeleteFile() {
	//Skip root folder
	if file.ParentId == uuid.Nil {
		return
	}
	DB.Unscoped().Delete(file)
}