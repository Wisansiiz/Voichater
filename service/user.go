package service

import (
	"Voichatter/dao"
	"Voichatter/models"
	"Voichatter/pkg/utils/jwt"
	"Voichatter/pkg/utils/translator"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func UserRegister(user *models.User) (err error) {
	if err = translator.ReErr(user); err != nil {
		return err
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
		return
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
	// 找到用户名为 "XXX" 的用户
	if err = dao.DB.Where("username = ?", user.Username).First(&user).Error; err != nil {
		return err
	}
	// 找到用户，获取他加入的服务器列表
	dao.DB.Table("server").
		Joins("JOIN member ON server.server_id = member.server_id").
		Where("member.user_id = ?", user.UserID).
		Find(&server)
	fmt.Printf("用户 %s 加入的服务器列表:\n", user.Username)
	fmt.Println(server)
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
