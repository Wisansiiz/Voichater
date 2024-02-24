package service

import (
	"Voichatter/dao"
	"Voichatter/models"
	"errors"
)

func GetServerMembers(serverId uint, users *[]models.UserList4Server) error {
	// 获取服务器的成员列表的用户信息
	var userIds []uint
	err := dao.DB.Model(&models.Member{}).
		Select("user_id").
		Where("server_id = ?", serverId).
		Find(&userIds).Error
	if err != nil {
		return errors.New("获取用户列表失败")
	}
	err = dao.DB.Model(&models.User{}).
		Select("user.user_id", "user.username", "user.email", "user.avatar_url", "member.s_permissions", "user.last_login_date").
		Joins("JOIN member ON user.user_id = member.user_id").
		Where("user.user_id IN (?) AND server_id = ?", userIds, serverId).
		Find(&users).Error
	if err != nil {
		return errors.New("获取用户列表失败")
	}
	return err
}
