package passwordless

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"github.com/guregu/null"
	"io"
	"net/url"
	"time"
)

const (
	TokenLength int           = 32
	TtlDuration time.Duration = 20 * time.Minute
)

type User struct {
	Id        int64       `db:"id"`
	Email     string      `db:"email"`
	Token     string      `db:"token"`
	Ttl       time.Time   `db:"ttl"`
	OriginUrl null.String `db:"originurl"`
}

// RefreshToken refreshes Ttl and Token for the User.
func (u *User) RefreshToken() error {
	token := make([]byte, TokenLength)
	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		return err
	}
	u.Token = base64.URLEncoding.EncodeToString(token)
	u.Ttl = time.Now().UTC().Add(TtlDuration)
	return nil
}

// IsValidToken returns a bool indicating that the User's current token hasn't
// expired and that the provided token is valid.
func (u *User) IsValidToken(token string) bool {
	if u.Ttl.Before(time.Now().UTC()) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(u.Token), []byte(token)) == 1
}

func (u *User) UpdateOriginUrl(originUrl *url.URL) error {
	var nsOrigin null.String
	if err := nsOrigin.Scan(originUrl.String()); err != nil {
		return err
	}

	u.OriginUrl = nsOrigin

	return nil
}
