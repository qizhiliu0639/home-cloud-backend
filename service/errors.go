package service

import "errors"

var (
	ErrRequestPara         = errors.New("invalid request")
	ErrInvalidOrPermission = errors.New("file or folder not exists or no permission")
	ErrDuplicate           = errors.New("duplicate file or folder")
	ErrSave                = errors.New("error when saving")
	ErrFoundFile           = errors.New("error when finding file")
	ErrConflict            = errors.New("file name conflicts with folder")
	ErrSystem              = errors.New("system error")
	ErrUsernameInvalid     = errors.New("username unavailable")
)
