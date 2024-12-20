package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var ErrGroupExists = errors.New("group already exists")
var ErrNoGroup = errors.New("group does not exist")
var ErrAlreadyInGroup = errors.New("user is already in group")
var ErrNotInGroup = errors.New("user is not in group")

type Group struct {
	Id   int64
	Name string
}

func GetGroup(db *sql.DB, gid int64) (*Group, error) {
	row := db.QueryRow(
		"SELECT * FROM grps WHERE id=?;",
		gid,
	)
	var group Group
	if err := row.Scan(
		&group.Id,
		&group.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoGroup
		} else {
			return nil, err
		}
	}

	return &group, nil
}

func FindGroup(db *sql.DB, name string) (*Group, error) {
	row := db.QueryRow(
		"SELECT * FROM grps WHERE name=?;",
		name,
	)
	var group Group
	if err := row.Scan(
		&group.Id,
		&group.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoGroup
		} else {
			return nil, err
		}
	}

	return &group, nil
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

func RemoveGroup(db *sql.DB, gid int64) error {
	// remove all users from the group

	result, err := db.Exec(
		"DELETE FROM user_group_rel WHERE gid=?;",
		gid,
	)
	if err != nil {
		return err
	}

	// remove group itself

	result, err = db.Exec(
		"DELETE FROM grps WHERE id=?;",
		gid,
	)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	switch {
	case err != nil:
		return err
	case n == 0:
		return ErrNoGroup
	}

	return nil
}

func GroupAddUser(db *sql.DB, uid int64, gid int64) error {
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

func GroupRemoveUser(db *sql.DB, uid int64, gid int64) error {
	result, err := db.Exec(
		"DELETE FROM user_group_rel WHERE uid=? AND gid=?;",
		uid, gid,
	)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return ErrNotInGroup
	}

	return nil
}

func IsUserInGroup(db *sql.DB, uid int64, gid int64) bool {
	row := db.QueryRow(
		"SELECT * FROM user_group_rel WHERE uid=? OR gid=?",
		uid, gid,
	)

	return row.Err() == nil
}

type ElementsToList int

const (
	GroupListUsers ElementsToList = 0
	UserListGroups ElementsToList = 1
)

func ListGroupsOrUsers(db *sql.DB, toList ElementsToList, id int64) ([]int64, error) {
	id1, id2 := [2]string{"uid", "gid"}[toList], [2]string{"gid", "uid"}[toList]
	rows, err := db.Query(
		fmt.Sprintf("SELECT %s FROM user_group_rel WHERE %s=?;", id1, id2),
		id,
	)
	if err != nil {
		return nil, err
	}

	var uids []int64
	for rows.Next() {
		var uid int64
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		uids = append(uids, uid)
	}

	return uids, nil
}
