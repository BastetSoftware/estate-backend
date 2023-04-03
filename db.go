package main

import (
	"fmt"
	"os"
	"time"

	"database/sql"
	"github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin = iota + 1
)

type UserInfo struct {
	Id         int64
	Login      string
	PassHash   []byte
	FirstName  string
	LastName   string
	Patronymic *string
	Role       int64
}

func (u UserInfo) Format() string {
	var patronymic string
	if u.Patronymic != nil {
		patronymic = *u.Patronymic
	} else {
		patronymic = "<none>"
	}
	return fmt.Sprint(
		u.Id, " ",
		u.Login, " ",
		u.PassHash, " ",
		u.FirstName, " ",
		u.LastName, " ",
		patronymic, " ",
		u.Role,
	)
}

func (u UserInfo) Register(db *sql.DB) error {
	var err error
	u.Id, err = registerUser(db, &u)
	return err
}

func registerUser(db *sql.DB, u *UserInfo) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO users (login, pass_hash, first_name, last_name, patronymic, role) VALUES (?,?,?,?,?,?)",
		u.Login, u.PassHash, u.FirstName, u.LastName, u.Patronymic, u.Role,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func getUserInfo(db *sql.DB, id int64) (*UserInfo, error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user UserInfo
	if err := row.Scan(
		&user.Id,
		&user.Login,
		&user.PassHash,
		&user.FirstName,
		&user.LastName,
		&user.Patronymic,
		&user.Role,
	); err != nil {
		return nil, err
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func findUserInfo(db *sql.DB, login string) (*UserInfo, error) {
	row := db.QueryRow("SELECT * FROM users WHERE login = ?", login)

	var user UserInfo
	if err := row.Scan(
		&user.Id,
		&user.Login,
		&user.PassHash,
		&user.FirstName,
		&user.LastName,
		&user.Patronymic,
		&user.Role,
	); err != nil {
		return nil, err
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

type Session struct {
	Id         int64
	ExpiryDate time.Time
	User       int64
}

func openSession(db *sql.DB, login string, pass string) (*Session, error) {
	user, err := findUserInfo(db, login)
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

func dbOpen() (*sql.DB, error) {
	var err error
	var db *sql.DB

	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "estate",
		AllowNativePasswords: true,
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
