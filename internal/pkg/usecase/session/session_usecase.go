package session

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"time"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

type JwtSessionUsecase struct {
	users domain.UserRepository
	cfg   *config.Session
}

func New(users domain.UserRepository, cfg *config.Session) domain.SessionUsecase {
	return &JwtSessionUsecase{
		users: users,
		cfg:   cfg,
	}
}

func (s JwtSessionUsecase) Login(creds domain.Credentials) (*domain.Session, error) {
	user, err := func() (*domain.User, error) {
		if creds.Email != "" {
			return s.users.GetByEmail(creds.Email)
		}
		return s.users.GetByName(creds.Name)
	}()
	if err != nil {
		return nil, err
	}

	if !checkPassword([]byte(creds.Password), user.Password) {
		return nil, app_errors.ErrInvalidCredentials
	}

	token, err := s.IssueToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.Session{
		ID:     token,
		UserID: user.ID,
	}, nil
}

func (s JwtSessionUsecase) IsLoggedIn(session *domain.Session) (bool, error) {
	token, err := jwt.Parse(session.ID, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})

	switch {
	case token != nil && token.Valid:
		subject, err := token.Claims.GetSubject()
		if err != nil {
			return false, err
		}

		session.UserID, err = utils.ParseUUID(subject)
		if err != nil {
			return false, err
		}
		return true, nil

	case errors.Is(err, jwt.ErrTokenMalformed):
		return false, app_errors.ErrInvalidCredentials

	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return false, app_errors.ErrInvalidSignature

	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return false, app_errors.ErrTokenExpired

	default:
		return false, err
	}
}

func (s JwtSessionUsecase) IssueToken(userID utils.UUID) (string, error) {
	issuedAt := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"iat": jwt.NewNumericDate(issuedAt),
		"exp": jwt.NewNumericDate(issuedAt.Add(time.Second * time.Duration(s.cfg.MaxAge))),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func checkPassword(got, expected []byte) bool {
	parts := bytes.Split(expected, []byte{':'})
	if len(parts) < 2 {
		return false
	}

	pass, salt := parts[0], parts[1]
	hashed := utils.Hash(got, salt)

	return slices.Equal(hashed, pass)
}
