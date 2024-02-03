package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"online-voice-channel/models"
	"time"
)

type MyCustomClaims struct {
	UserID       uint   `json:"userID"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
	GrantScope   string `json:"grant_scope"`
	jwt.RegisteredClaims
}

// 签名密钥
const signKey = "my jwt"

func GenerateToken(user models.User) (string, error) {
	claim := MyCustomClaims{
		UserID:       user.UserID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		GrantScope:   "read_user_info",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "Auth_Server", // 签发者
			Subject: user.Username, // 签发对象
			//Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
			//ID:        "123",                                           // wt ID, 类似于盐值
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(signKey))
	return token, err
}
func ParseToken(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(signKey), nil //返回签名密钥
		})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("claim invalid")
	}
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}
	return claims, nil
}
