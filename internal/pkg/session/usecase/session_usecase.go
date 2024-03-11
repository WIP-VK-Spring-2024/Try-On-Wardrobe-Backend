package usecase

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
	"slices"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type SessionUsecase struct {
	users    domain.UserRepository
	sessions domain.SessionRepository
}

func NewSessionUsecase(users domain.UserRepository, sessions domain.SessionRepository) domain.SessionUsecase {
	return &SessionUsecase{
		users:    users,
		sessions: sessions,
	}
}

func (s SessionUsecase) Login(creds domain.Credentials) (*domain.Session, error) {
	user, err := s.users.GetByName(creds.Name)
	if err != nil {
		return nil, err
	}

	if !checkPassword([]byte(creds.Password), user.Password) {
		return nil, app_errors.ErrInvalidCredentials
	}
	return nil, nil
}

func (s SessionUsecase) Logout(sessionID string) error {
	return s.sessions.Delete(sessionID)
}

func (s SessionUsecase) Register(user *domain.User) (*domain.Session, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, app_errors.New(err)
	}

	user.Password = slices.Concat(hash(user.Password, salt), []byte{':'}, salt)
	err = s.users.Create(user)
	if err != nil {
		return nil, app_errors.New(err)
	}

	session := domain.Session{
		ID:     uuid.NewString(),
		UserID: user.ID,
	}

	err = s.sessions.Put(session)
	if err != nil {
		return nil, errors.Join(err, app_errors.ErrSessionNotInitialized)
	}

	return &session, nil
}

func checkPassword(got, expected []byte) bool {
	parts := bytes.Split(expected, []byte{':'})
	if len(parts) < 2 {
		return false
	}

	pass, salt := parts[0], parts[1]
	hashed := hash(got, salt)

	return slices.Equal(hashed, pass)
}

func hash(pass, salt []byte) []byte {
	bytes := argon2.IDKey(pass, salt, 1, 64*1024, 4, 32)
	result := make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
	base64.StdEncoding.Encode(result, bytes)
	return result
}

var intMax *big.Int = big.NewInt(64 * 1024)

func generateSalt() ([]byte, error) {
	salt, err := rand.Int(rand.Reader, intMax)
	if err != nil {
		return nil, err
	}
	return salt.Bytes(), nil
}
