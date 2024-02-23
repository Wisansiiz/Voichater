package service

import (
	"Voichatter/dao"
	"Voichatter/models"
	"Voichatter/pkg/utils/jwt"
	"Voichatter/pkg/utils/translator"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func UserRegister(user *models.User) (err error) {
	if err = translator.ReErr(user); err != nil {
		log.Println(err)
		return errors.New("注册信息未填写正确")
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return err
	}
	user.PasswordHash = string(pwd)
	user.RegistrationDate = time.Now()
	err = dao.DB.Create(&user).Error
	return err
}

func UserLogin(user *models.UserLoginResponse) (token string, err error) {
	if err = translator.ReErr(user); err != nil {
		log.Println(err)
		return "", errors.New("登录信息未填写正确")
	}
	// 数据库内的用户名密码
	var u models.User
	if err = dao.DB.Where("username", user.Username).Find(&u).Error; err != nil {
		return
	}
	pwd := user.PasswordHash
	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pwd)); err != nil {
		return
	}
	token, err = jwt.GenerateToken(u)
	if err != nil {
		return
	}
	if err = dao.DB.Model(&u).Update("last_login_date", time.Now()).Error; err != nil {
		return
	}
	return token, err
}

func FindUserServersList(user *models.User, server *[]models.Server) (err error) {
	// 找到用户，获取他加入的服务器列表
	dao.DB.Table("server").
		Joins("JOIN member ON server.server_id = member.server_id").
		Where("member.user_id = ?", user.UserID).
		Find(&server)
	return err
}

func UserLogout(token string) (err error) {
	if _, err = dao.RedisClient.Incr(dao.RedisContext, token).Result(); err != nil {
		return err
	}
	if _, err = dao.RedisClient.Expire(dao.RedisContext, token, time.Hour).Result(); err != nil {
		return err
	}
	return err
}

func CreateServer(user *models.User, server *models.Server) (err error) {
	// 创建服务器，创建者id为user.id，并且权限为admin
	s := models.Server{
		ServerName:    server.ServerName,
		CreatorUserID: user.UserID,
		CreateDate:    time.Now(),
		ServerType:    server.ServerType,
	}
	if err = translator.ReErr(s); err != nil {
		log.Println(err)
		return errors.New("服务器信息未填写正确")
	}
	if err = dao.DB.Create(&s).Error; err != nil {
		return errors.New("创建服务器失败")
	}
	m := models.Member{
		UserID:       s.CreatorUserID,
		ServerID:     s.ServerID,
		JoinDate:     time.Now(),
		SPermissions: "admin",
		CPermissions: "admin",
	}
	if err = dao.DB.Create(&m).Error; err != nil {
		return errors.New("创建服务器成员失败")
	}
	return err
}

func JoinServer(member *models.Member) error {
	err := dao.DB.Model(&models.Server{}).
		Select("server_id").
		Where("server_id = ? AND server_type = 'public' ", member.ServerID).
		Find(&member).Error
	if err != nil {
		return errors.New("服务器不存在")
	}
	// 判断是否已经加入过服务器
	err = dao.DB.Model(&models.Member{}).
		Select("member_id").
		Where("server_id = ? AND user_id = ? ", member.ServerID, member.UserID).
		Find(&member).Error
	if err != nil || member.MemberID != 0 {
		return errors.New("已经加入过服务器")
	}
	m := models.Member{
		ServerID:     member.ServerID,
		UserID:       member.UserID,
		JoinDate:     time.Now(),
		SPermissions: "member",
		CPermissions: "member",
	}
	// 添加成员
	err = dao.DB.Model(&models.Member{}).Create(&m).Error
	if err != nil {
		return errors.New("添加成员失败")
	}
	return err
}
