package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"online-voice-channel/dao"
	"online-voice-channel/models"
)

func UserRegister(user *models.User) (err error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return err
	}
	user.PasswordHash = string(pwd)
	err = dao.DB.Create(&user).Error
	return err
}

func FindUserServersList(user *models.User, servers *[]models.Server) (err error) {
	// 找到用户名为 "XXX" 的用户
	if err = dao.DB.Where("username = ?", user.Username).First(&user).Error; err != nil {
		return err
	}

	// 找到用户，获取他加入的服务器列表
	dao.DB.Table("servers").
		Joins("JOIN members ON servers.server_id = members.server_id").
		Where("members.user_id = ?", user.UserID).
		Find(&servers)
	fmt.Printf("用户 %s 加入的服务器列表:\n", user.Username)
	return err
}