package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/xxandjg/ginbase/global"
	"time"
)

type JWT struct {
}

type MyClaims struct {
	UserId             string `json:"userId"`
	jwt.StandardClaims        // 标准Claims结构体，可设置8个标准字段
}

const PriKey = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIC3TBXBgkqhkiG9w0BBQ0wSjApBgkqhkiG9w0BBQwwHAQIVq07cH81c9wCAggA
MAwGCCqGSIb3DQIJBQAwHQYJYIZIAWUDBAECBBDClNpP6FIw76TjU4NlOuCfBIIC
gErFYrdu8B0RTZwSrBduqqAoTeY3elQO5LCBN80PsQ6K1EnU2cmQz1w8cqJfS7eB
vI/LHJ7AqDnmBs84wuISgLYN4x7uZNgd0Q5f7cfYS5RcUXHRVD7/mI6YOK4YfzR3
Q+fhMcyPpeKlXvzLzR7TIqBpoZw/NYTdd6KEkUKB1LPpQH5prjh7lUmDDyXwgfRu
VFhDluvLD/X9feI+Gznmds5rP61QHSFfmD5qASHMPcdtJsyO2DRTH9xP3oPC809F
oMsAAEfRKNKOlEMTStn5acne3krjS9jdyVWmG0IyHb18AG/u/zWyk33T7rrrX7J4
TlWKX0ZfPevIMhhGs8h1ZQQKmNpqbU8fO3mNfwFtKvf1dPNhhkQ4X4MM+rLFwkCp
ENxG3E1hI93QKLA1AJSgZ8YcUdqhsVkBlHi6vmZzAfrkt9Y2Ax/+UIeSGX5ms0fT
jzWYVXNwcL8Vl3/jffd0KEvoZi+ct+Hjb9SmC0I/ju9pTiREb7iAJGpQjXY/GbRL
BWxqOaajhtmvumiQo6nArviQq+OyF7jtOfnLpgAg47Ht8wbWSOZSosOhvbs+Jt9P
Cm6tQAf6x7hj3twrluZ9zPUyL5o8fkGDkddRC4Nzx6k0QCwbtNRabkO42t/qDPsB
3uoj9D3gHyFGqo+4p9pbdSxVNA57Yesp7ZIUrMyii9m++bGlL3swVeuyXLyX9qB5
JHmDIjqEl7NdzmcboVOmHTsbqqI65ZGYlXBZfz2AbsXR4c6BPSSS+BrIw5OR2tAB
Qh5PZUPNaYVh5reFi5343bUMauU1WfcqYZTyOAD+CAgNkhCxT9VRMHZZBAsS3avz
R13vXANWEFEOFsoHqM8Qc64=
-----END ENCRYPTED PRIVATE KEY-----
`

const PubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGI2Nitvz9tW1yyWdPiUVPUJzo
mPbTNjhWcn0Mu+iky64RzJWYFLUZ8oHxKPmyt1FZmX1jQqefzW8MEOku2bfWBeUj
1Km4jl1xXOT9YQoYDm/XkmqawTZs4SBp6dUkURon9mbJI3t/6Zosx1kwMdYPlscb
e6zSkntezHntM1FCrQIDAQAB
-----END PUBLIC KEY-----
`

const TokenExpireDuration = time.Hour * 24 * 7

const MySecret = "$2a$10$ZJHCn61zJLYeVVkwFl/aRuvfLrQVCCvb6Mrc.9zP.gnzq6xKvEIhq"

// GenerateToken 传入雪花ID，生成签名的密钥
func (j *JWT) GenerateToken(userId string) (string, error) {
	expirationTime := time.Now().Add(TokenExpireDuration)
	claims := &MyClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    MySecret,
		},
	}
	// 生成Token，指定签名算法和claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名
	if tokenString, err := token.SignedString([]byte(PriKey)); err != nil {
		return "", err
	} else {
		return tokenString, nil
	}

}

// ParseToken 通过jwt.ParseWithClaims返回的Token结构体取出Claims结构体
func (j *JWT) ParseToken(t string) (*MyClaims, global.Error) {
	//zap.L().Info(t)
	token, err := jwt.ParseWithClaims(t, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(PriKey), nil
	})
	if err != nil {
		return nil, global.FAIL
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, global.SUCCESS
	}
	return nil, global.FAIL
}
