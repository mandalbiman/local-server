package svc

import (
	"errors"

	"bats.com/local-server/api/models"
)

type Auth interface {
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
}

type AuthSvcImpl struct {
	//db db.Database
}

func NewAuth() Auth {
	return &AuthSvcImpl{}
}

func (a AuthSvcImpl) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	//user, err := a.db.GetUser(req.Name)
	if "hi-poltu" != req.Password || "poltu" != req.Name {
		return nil, errors.New("invalid credentials")
	}

	id := "poltu"
	roles := []string{"admin"}

	jwtToken, _ := GenerateJWT(id, roles)
	return &models.LoginResponse{
		Roles: &roles,
		Token: &jwtToken,
		User:  &id,
	}, nil

}
