package store

import (
	"database/sql"
	"time"

	"github.com/Nathac/go-api/internals/store/tokens"
)

type postgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *postgresTokenStore {
	return &postgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(*tokens.Token) error
	CreateToken(int, time.Duration, string) (*tokens.Token, error)
	DeleteToken(int) error
}

func (t *postgresTokenStore) CreateToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *postgresTokenStore) Insert(token *tokens.Token) error {
	query := `INSERT INTO tokens(hash,user_id, expiry, scope) VALUES($1,$2,$3,$4)`

	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	if err != nil {
		return err
	}
	return nil
}

func (t postgresTokenStore) DeleteToken(UserID int) error {
	query := `DELETE token where user_id=$1`

	result, err := t.db.Exec(query, UserID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
