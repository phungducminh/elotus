package auth

import (
	"fmt"
	"strconv"
	"time"

	"elotus.com/hackathon/storage"
	"elotus.com/hackathon/storage/query"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = fmt.Errorf("auth: invalid credentials")
	ErrInvalidOrExpired   = fmt.Errorf("auth: invalid token or expired")
	ErrUsernameNotUnique  = fmt.Errorf("auth: username not unique")
)

type Auth interface {
	Register(req *RegisterRequest) (*RegisterResponse, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	Verify(req *VerifyRequest) (*VerifyResponse, error)
}

type auth struct {
	storage          storage.Storage
	secretKey        []byte
	expiresInSeconds int
	lg               *zap.Logger
}

func NewAuth(lg *zap.Logger, storage storage.Storage, secretKey []byte, expiresInSeconds int) Auth {
	return &auth{
		lg:               lg,
		storage:          storage,
		secretKey:        secretKey,
		expiresInSeconds: expiresInSeconds,
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

	expiresAt := time.Now().Add(time.Second * time.Duration(au.expiresInSeconds)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expiresAt,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "elotus",
		Subject:   strconv.FormatInt(user.ID, 10),
	})

	accessToken, err := token.SignedString(au.secretKey)
	if err != nil {
		return nil, err
	}

	resp := &LoginResponse{
		AccessToken: accessToken,
		ExpiresAt:   strconv.FormatInt(expiresAt, 10),
		UserId:      strconv.FormatInt(user.ID, 10),
	}

	return resp, nil
}

func (au *auth) Verify(req *VerifyRequest) (*VerifyResponse, error) {
	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(req.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return au.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	resp := &VerifyResponse{
		UserId: claims.Subject,
	}
	return resp, nil
}
