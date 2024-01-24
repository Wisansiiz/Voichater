package models

import (
	"gorm.io/gorm"
	"online-voice-channel/pkg/utils/jwt"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	*gorm.Model
	//Remember bool   `json:"remember"`
	//Code     string `json:"code"`
	//UUID     string `json:"uuid"`
	//Role     string `json:"role"`
}

func LoginOn(login *Login) (err error) {
	if login.Username == "123" && login.Password == "123" {
		if login.Token, err = jwt.GenerateToken(login.ID, login.Username); err == nil {
			return
		}
	}
	return
}
