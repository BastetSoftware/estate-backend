package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

var ErrUserExists = errors.New("user already exists")
var ErrNoUser = errors.New("user does not exist")

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
	u.Id, err = RegisterUser(db, &u)
	return err
}

func RegisterUser(db *sql.DB, u *UserInfo) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO users (login, pass_hash, first_name, last_name, patronymic, role) VALUES (?,?,?,?,?,?)",
		u.Login, u.PassHash, u.FirstName, u.LastName, u.Patronymic, u.Role,
	)
	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == 1062 {
				return 0, ErrUserExists
			}
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetUserInfo(db *sql.DB, id int64) (*UserInfo, error) {
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
		if err == sql.ErrNoRows {
			return nil, ErrNoUser
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func FindUserInfo(db *sql.DB, login string) (*UserInfo, error) {
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
		if err == sql.ErrNoRows {
			return nil, ErrNoUser
		} else {
			return nil, err
		}
	}

	return &user, nil
}
