package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"home-cloud/models"
	"home-cloud/service"
	"net/http"
	"strconv"
)

// GetUserList get user list in the system
func GetUserList(c *gin.Context) {
	user := c.Value("user").(*models.User)
	users, err := service.GetUserList(user)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
		return
	}
	resUserInfo := make([]gin.H, len(users))
	for i, v := range users {
		resUserInfo[i] = gin.H{
			"Name":       v.Username,
			"Status":     v.Status,
			"Quota":      v.Storage,
			"Encryption": v.Encryption != 0,
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": 0, "users": resUserInfo})
}

// DeleteUser delete a user
func DeleteUser(c *gin.Context) {
	user := c.Value("user").(*models.User)
	deleteUser := c.PostForm("delete_user")
	if deleteUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Please input username"})
		return
	}
	if user.Username == deleteUser {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "You cannot delete yourself"})
		return
	}
	err := service.DeleteUser(user, deleteUser)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}

// SetUserQuota set user storage quota
func SetUserQuota(c *gin.Context) {
	modifiedUser := c.PostForm("modified_user")
	if modifiedUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Please input username"})
		return
	}
	size := c.PostForm("quota")
	if size == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Please input new size"})
		return
	}
	newSize, err := strconv.ParseUint(size, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Please input valid new size"})
		return
	}
	err = service.SetUserQuota(modifiedUser, newSize)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}

// ToggleAdmin set user permission
func ToggleAdmin(c *gin.Context) {
	user := c.Value("user").(*models.User)
	toggleUser := c.PostForm("toggle_user")
	if toggleUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Please input username"})
		return
	}
	if user.Username == toggleUser {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "You cannot change your own permission! "})
		return
	}
	err := service.ToggleAdmin(user, toggleUser)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}

// ResetUserPassword reset password for a user
func ResetUserPassword(c *gin.Context) {
	resetUser := c.PostForm("reset_user")
	if resetUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Please input username"})
		return
	}
	res, err := service.ResetUserPassword(resetUser)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0, "result": res})
	}
}
