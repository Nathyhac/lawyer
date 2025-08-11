package store

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int      `json:"id"`
	First_name    string   `json:"first_name"`
	Last_name     string   `json:"last_name"`
	UserName      string   `json:"username"`
	Email         string   `json:"email"`
	Phone_number  string   `json:"phone_number"`
	Hash_Password Password `json:"hash_password"`
}
type Password struct {
	Plaintext *string
	Hash      []byte
}

func (p *Password) Set(plaintText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintText), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintText
	p.Hash = hash

	return nil
}

func (p *Password) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, err
		default:
			return false, err
		}

	}

	return true, nil
}

type userPostgresDb struct {
	db *sql.DB
}

func NewUserpostgresDB(db *sql.DB) *userPostgresDb {
	return &userPostgresDb{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) (*User, error)

	GetUserUsername(userName string) (*User, error)
}

func (pg *userPostgresDb) CreateUser(user *User) (*User, error) {

	query := `INSERT INTO users(first_name,last_name, username,email, phone_number,hash_password) 
	VALUES($1,$2,$3,$4,$5,$6) RETURNING id`

	err := pg.db.QueryRow(query, user.First_name, user.Last_name, user.UserName, user.Email, user.Phone_number, user.Hash_Password.Hash).Scan(&user.ID)
	if err != nil {
		fmt.Printf("error creating a user: %v", err)
		return nil, err
	}

	return user, nil

}

func (pg *userPostgresDb) GetUserUsername(username string) (*User, error) {

	user := &User{
		Hash_Password: Password{},
	}
	query := `SELECT id, first_name, last_name,username, email,phone_number, hash_password FROM users WHERE username=$1`
	err := pg.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.First_name,
		&user.Last_name,
		&user.UserName,
		&user.Email,
		&user.Phone_number,
		&user.Hash_Password.Hash,
	)
	if err != nil {
		return nil, fmt.Errorf("error while querying with username:%v", err)
	}
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no rows affected: %v", err)
	}

	return user, nil

}

func (pg *userPostgresDb) UpdateUser(user *User) error {

	query := `UPDATE Users WHERE first_name=$1, last_name=$2,email=$3,username=$4, phone_number=$5`

	result, err := pg.db.Exec(query, user)
	if err != nil {
		return fmt.Errorf("error executing the query:%v", err)
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
