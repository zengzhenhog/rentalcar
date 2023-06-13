package token

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTTokenGen struct {
	privateKey *rsa.PrivateKey
	issuer     string
	nowFunc    func() time.Time
}

func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		issuer:     issuer,
		nowFunc:    time.Now,
		privateKey: privateKey,
	}
}

func (t *JWTTokenGen) GenerateToken(accountID string, expire time.Duration) (string, error) {
	nowSec := t.nowFunc().Unix()
	// 基础配置jwt.StandardClaims， MapClaims可以加自定义字段，注意tkn结果的字段是按照字母序来排
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    t.issuer,
		IssuedAt:  nowSec,
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   accountID,
	})
	return tkn.SignedString(t.privateKey) // 返回签名后的结果
}
