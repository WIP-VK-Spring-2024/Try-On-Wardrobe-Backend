package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"try-on/internal/pkg/api_errors"
	"try-on/internal/pkg/session"
)

var ErrBadToken = errors.New("invalid token format")

type CsrfTokenManager struct {
	secret []byte
}

func New(secret string) *CsrfTokenManager {
	return &CsrfTokenManager{secret: []byte(secret)}
}

func (token *CsrfTokenManager) Create(s *session.Session, expiration int64) string {
	h := hmac.New(sha256.New, token.secret)
	data := fmt.Sprintf("%s:%d:%d", s.ID, s.UserID, expiration)
	h.Write([]byte(data))
	result := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return result + ":" + strconv.FormatInt(expiration, 10)
}

func (token *CsrfTokenManager) Check(s *session.Session, csrfToken string) (bool, error) {
	if csrfToken == "" || s == nil {
		return false, nil
	}

	tokenData := strings.Split(csrfToken, ":")
	if len(tokenData) != 2 {
		return false, api_errors.ErrInvalidToken
	}

	expiration, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, api_errors.ErrInvalidToken
	}

	if expiration < time.Now().Unix() {
		return false, api_errors.ErrExpired
	}

	h := hmac.New(sha256.New, token.secret)
	data := fmt.Sprintf("%s:%d:%d", s.ID, s.UserID, expiration)
	h.Write([]byte(data))

	expectedMAC := h.Sum(nil)
	gotMAC, err := base64.StdEncoding.DecodeString(tokenData[0])
	if err != nil {
		return false, api_errors.ErrInvalidToken
	}

	return hmac.Equal(gotMAC, expectedMAC), nil
}
