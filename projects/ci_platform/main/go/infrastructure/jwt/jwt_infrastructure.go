package infraJWT

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg"

	infraConfig "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/infrastructure/config"
	useJWT "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/jwt"
)

type interactor struct {
	conf   *infraConfig.JWT
	secret []byte
}

func New(conf *infraConfig.JWT) useJWT.Interactor {
	return &interactor{conf: conf, secret: conf.GetSecret()}
}

func (ji *interactor) New(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"admin": user.Admin,
		"exp": time.Now().Add(
			time.Hour * time.Duration(ji.conf.Expiration),
		).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ji.secret)
}

func (ji *interactor) Claims(token interface{}) (*entity.JWTClaims, error) {
	claimErr := func(format string, a ...any) error {
		return fmt.Errorf("cannot get claims: "+format, a...)
	}
	if token == nil {
		return nil, claimErr("jwt is nil")
	}
	tokenAsserted, ok := token.(*jwt.Token)
	if !ok {
		return nil, claimErr("invalid token: %+v", token)
	}
	claimsMap, ok := tokenAsserted.Claims.(jwt.MapClaims)
	if !ok {
		return nil, claimErr("invalid claim map: %+v", tokenAsserted.Claims)
	}
	claims := &entity.JWTClaims{}
	if err := mapstructure.Decode(token, claimsMap); err != nil {
		return nil, claimErr("%w", err)
	}
	if err := pkg.ValidateStruct(claims); err != nil {
		return nil, claimErr("%w", err)
	}
	return claims, nil
}
