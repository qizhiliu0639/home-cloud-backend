package controllers

import "home-cloud/service"

func GetErrorMessage(err error) (res string) {
	switch err {
	case service.ErrRequestPara:
		res = "Invalid Request"
	case service.ErrDuplicate:
		res = "Duplicate File or Folder Name"
	case service.ErrInvalidOrPermission:
		res = "Invalid File or Folder or Permission Denied"
	case service.ErrSave:
		res = "Errors in Saving or Creating"
	case service.ErrFoundFile:
		res = "Errors in Finding File or Folder"
	case service.ErrConflict:
		res = "Conflict in File Name or Folder Name"
	case service.ErrSystem:
		res = "System Error"
	}
	return
}
