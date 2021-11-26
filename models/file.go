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
	ID uuid.UUID `gorm:"type:char(36);primaryKey"`
	// IsDir 1 for folder and 0 for file
	IsDir int `gorm:"default:0;not null"`
	// Name The name of the file or folder. For the root folder, it will be the username
	Name string `gorm:"type:varchar(191);not null;uniqueIndex:idx_only_one"`
	// ParentId uuid.Nil for root folder
	ParentId  uuid.UUID `gorm:"type:char(36);not null;uniqueIndex:idx_only_one"`
	OwnerId   uuid.UUID `gorm:"type:char(36);not null"`
	CreatorId uuid.UUID `gorm:"type:char(36);not null"`
	Size      uint64    `gorm:"default:0;not null"`
	FileType  string    `gorm:"default:'other'"`
	RealPath  string    `gorm:"not null"`
	Favorite  int       `gorm:"default:0"`

	// Position The position of file. This field will be ignored in the database
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
	err = DB.Where(&File{ParentId: file.ID, OwnerId: file.OwnerId}).Order("is_dir desc").Order("name").Find(&children).Error
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

func (file *File) AddFavorite() error {
	return DB.Model(&file).Update("favorite", 1).Error
}

func (file *File) CancelFavorite() error {
	return DB.Model(&file).Update("favorite", 0).Error
}
