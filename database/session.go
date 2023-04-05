package database

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type Session struct {
	Id         int64
	Token      []byte
	ExpiryDate time.Time
	User       int64
}

const tokenLength = 32
const tokenAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func OpenSession(db *sql.DB, login string, pass string) (*Session, error) {
	user, err := FindUserInfo(db, login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass))
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	// generate a token
	token := make([]byte, tokenLength)
	for i := range token {
		token[i] = tokenAlphabet[rand.Intn(len(tokenAlphabet))]
	}

	var session = Session{
		Id:         0,
		Token:      token,
		ExpiryDate: time.Now().AddDate(0, 0, 2),
		User:       user.Id,
	}
	result, err := db.Exec(
		"INSERT INTO sessions (token, expiry_date, user) VALUES (?,?,?)",
		session.Token,
		session.ExpiryDate,
		session.User,
	)
	if err != nil {
		return nil, err
	}

	session.Id, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func CloseSession(db *sql.DB, token []byte) error {
	result, err := db.Exec(
		"DELETE FROM sessions WHERE token=?",
		token,
	)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n < 1 {
		return fmt.Errorf("no rows were affected")
	}

	return nil
}
