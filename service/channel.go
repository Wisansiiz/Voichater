package service

import (
	"Voichatter/dao"
	"Voichatter/models"
	"Voichatter/pkg/utils/translator"
	"errors"
	"log"
	"time"
)

func CreateChannel(channel *models.Channel, userId uint) error {
	// 查询用户是否在服务器中
	var member models.Member
	err := dao.DB.Model(&models.Member{}).
		Select("member_id").
		Where("server_id = ? AND user_id = ?", channel.ServerID, userId).
		Where("s_permissions = ?", "admin").
		Or("s_permissions = ?", "super_admin").
		Find(&member).Error
	if err != nil {
		return errors.New("查询出错")
	}
	if member.MemberID == 0 {
		return errors.New("用户不在服务器中或权限不足")
	}
	// 创建频道
	var c = &models.Channel{
		ServerID:     channel.ServerID,
		ChannelName:  channel.ChannelName,
		CreateUserId: userId,
		Type:         channel.Type,
		CreationDate: time.Now(),
	}
	if err = translator.ReErr(c); err != nil {
		log.Println(err)
		return errors.New("频道信息未填写正确")
	}
	err = dao.DB.Create(&c).Error
	if err != nil {
		return errors.New("创建频道失败")
	}
	return err
}
