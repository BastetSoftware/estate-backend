package database

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Session struct {
	Id         int64
	ExpiryDate time.Time
	User       int64
}

func OpenSession(db *sql.DB, login string, pass string) (*Session, error) {
	user, err := FindUserInfo(db, login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass))
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	var session = Session{
		Id:         0,
		ExpiryDate: time.Now().AddDate(0, 0, 2),
		User:       user.Id,
	}
	result, err := db.Exec(
		"INSERT INTO sessions (expiry_date, user) VALUES (?,?)",
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
