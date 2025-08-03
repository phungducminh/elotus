package auth

import (
	"fmt"
	"strconv"
	"time"

	"elotus.com/hackathon/storage"
	"elotus.com/hackathon/storage/query"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = fmt.Errorf("auth: invalid credentials")
	ErrUsernameNotUnique  = fmt.Errorf("auth: username not unique")
)

type Auth interface {
	Register(req *RegisterRequest) (*RegisterResponse, error)
	Login(req *LoginRequest) (*LoginResponse, error)
}

type auth struct {
	storage   storage.Storage
	secretKey []byte
}

func NewAuth(storage storage.Storage, secretKey []byte) Auth {
	return &auth{
		storage:   storage,
		secretKey: secretKey,
	}
}

func (au *auth) Register(req *RegisterRequest) (*RegisterResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// TODO: @check in cache
	_, err = au.storage.GetUserByUserName(req.Username)
	if err != storage.ErrNotFound {
		return nil, ErrUsernameNotUnique
	}

	user := &query.InsertUserParams{
		Username:       req.Username,
		HashedPassword: string(hashed),
	}
	id, err := au.storage.InsertUser(user)
	if err != nil {
		return nil, err
	}

	resp := &RegisterResponse{
		UserId: strconv.FormatInt(id, 10),
	}
	return resp, nil
}

func (au *auth) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := au.storage.GetUserByUserName(req.Username)
	if err == storage.ErrNotFound {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "elotus",
		Subject:   user.Username,
	})

	accessToken, err := t.SignedString(au.secretKey)
	if err != nil {
		return nil, err
	}

	resp := &LoginResponse{
		AccessToken: accessToken,
		UserId:      strconv.FormatInt(user.ID, 10),
	}

	return resp, nil
}
