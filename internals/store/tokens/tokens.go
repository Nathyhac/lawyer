package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)
const(
	ScopeAuth = "authentication"
)

type Token struct {
	PlaintText string    `json:"token"`
	Hash       []byte    `json:"-"`
	UserID     int       `json:"-"`
	Expiry     time.Time `json:"expiry"`
	Scope      string    `json:"-"`
}

func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	emptyByte := make([]byte, 32)
	_, err := rand.Read(emptyByte)
	if err != nil {
		return nil, err
	}

	token.PlaintText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyByte)
	hash := sha256.Sum256([]byte(token.PlaintText))
	token.Hash = hash[:]
	return token, nil
}
