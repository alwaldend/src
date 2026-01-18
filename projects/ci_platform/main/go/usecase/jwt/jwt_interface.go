package useJWT

import "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"

type Interactor interface {
	// generate user jwt
	New(user *entity.User) (string, error)
	// get claims from jwt
	Claims(token interface{}) (*entity.JWTClaims, error)
}
