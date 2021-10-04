package models

import (
	"errors"
	"gorm.io/gorm"
)

type File struct {
	// 表字段
	gorm.Model
	ID uint64 `gorm:"primaryKey;AUTO_INCREMENT"`
	//1 for folder and 0 for file
	IsDir int    `gorm:"default:0;not null"`
	Name  string `gorm:"type:varchar(255);uniqueIndex:idx_only_one"`
	//null for root folder
	ParentId  uint64 `gorm:"default:null;uniqueIndex:idx_only_one"`
	OwnerId   uint64
	CreatorId uint64
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

func (file *File) CreateFile() error {
	return DB.Create(file).Error
}

func GetFileByID(fid uint64) (*File, error) {
	var file File
	err := DB.Where(&File{ID: fid}).First(&file).Error
	return &file, err
}
func GetFileByName(name string, owner *User) (*File, error) {
	var file File
	err := DB.Where(&File{Name: name, OwnerId: owner.ID}).First(&file).Error
	return &file, err
}

func (file *File) GetChildInFolder() ([]*File, error) {
	if file.IsDir == 0 {
		return []*File{}, errors.New("not a folder")
	}
	var child []*File
	err := DB.Where(&File{ParentId: file.ID}).Find(&child).Error
	return child, err
}

func NewFile() *File {
	return &File{}
}
