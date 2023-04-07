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
	Patronymic string
	Role       int64
}

func (u UserInfo) Format() string {
	var patronymic = u.Patronymic
	if patronymic == "" {
		patronymic = "-"
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
		"INSERT INTO users (login, pass_hash, first_name, last_name, patronymic, role) VALUES (?,?,?,?,?,?);",
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
	row := db.QueryRow("SELECT * FROM users WHERE id = ?;", id)

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
	row := db.QueryRow("SELECT * FROM users WHERE login = ?;", login)

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

func UserChangeLogin(db *sql.DB, id int64, newLogin string) error {
	_, err := db.Exec("UPDATE users SET login=? WHERE id=?;", newLogin, id)
	switch e := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		if e.Number == 1062 {
			return ErrUserExists
		}
	default:
		if err == sql.ErrNoRows {
			return ErrNoUser
		} else {
			return err
		}
	}

	return nil
}

func UserChangePasswordHash(db *sql.DB, id int64, pass_hash []byte) error {
	_, err := db.Exec("UPDATE users SET pass_hash=? WHERE id=?;", pass_hash, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoUser
		} else {
			return err
		}
	}

	return nil
}

func UserChangeName(db *sql.DB, id int64, nameType int, newName string) error {
	var nameTypes = [2]string{"first_name", "last_name"}
	if nameType >= len(nameTypes) || nameType < 0 {
		return fmt.Errorf("invalid name type")
	}

	_, err := db.Exec("UPDATE users SET "+nameTypes[nameType]+"=? WHERE id=?;", newName, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoUser
		} else {
			return err
		}
	}

	return nil
}
