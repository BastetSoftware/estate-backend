package database

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
)

var ErrGroupExists = errors.New("group already exists")
var ErrAlreadyInGroup = errors.New("user is already in group")

type Group struct {
	Id   int64
	Name string
}

func CreateGroup(db *sql.DB, name string) (*Group, error) {
	result, err := db.Exec(
		"INSERT INTO grps (name) VALUES (?);",
		name,
	)
	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == 1062 {
				return nil, ErrGroupExists
			}
		}
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Group{Id: id, Name: name}, nil
}

func AddUserToGroup(db *sql.DB, uid int64, gid int64) error {
	_, err := db.Exec(
		"INSERT INTO user_group_rel (uid, gid) VALUES (?,?);",
		uid, gid,
	)
	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == 1062 {
				return ErrAlreadyInGroup
			}
		}
		return err
	}

	return nil
}

func IsUserInGroup(db *sql.DB, uid int64, gid int64) bool {
	row := db.QueryRow(
		"SELECT * FROM user_group_rel WHERE uid=? OR (gid=?)",
		uid, gid,
	)

	return row.Err() == nil
}
