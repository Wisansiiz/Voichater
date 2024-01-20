package models

import "online-voice-channel/test"

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	//Remember bool   `json:"remember"`
	//Code     string `json:"code"`
	//UUID     string `json:"uuid"`
	//Role     string `json:"role"`
}

func LoginOn(login *Login) (err error) {
	if login.Username == "123" && login.Password == "123" {
		if login.Token, err = test.GenerateToken(login.Username); err == nil {
			return
		}
	}
	return
}
