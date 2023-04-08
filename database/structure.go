package database

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

var ErrStructExists = errors.New("struct already exists")
var ErrNoStruct = errors.New("struct does not exist")

type StructInfo struct {
	Id          int64
	Name        string
	Description string
	District    string
	Region      string
	Address     string
	Type        string
	State       string
	Area        int32
	Owner       string
	Actual_user string
	Gid         int64
	Permissions int8
}

func (strct *StructInfo) AddStruct(db *sql.DB) error {
	result, err := db.Exec(
		"INSERT INTO objects (name, description, district, region, address, type, state, area, owner, actual_user, gid, permissions) VALUES (?,?,?,?,?,?,?,?,?,?,?,?);",
		strct.Name, strct.Description, strct.District, strct.Region,
		strct.Address, strct.Type, strct.State, strct.Area,
		strct.Owner, strct.Actual_user, strct.Gid,
		strct.Permissions,
	)
	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == 1062 {
				return ErrStructExists
			}
		}
		return err
	}

	id, err := result.LastInsertId()
	strct.Id = id
	if err != nil {
		return err
	}

	return nil
}

func GetStructInfo(db *sql.DB, id int64) (*StructInfo, error) {
	row := db.QueryRow("SELECT * FROM objects WHERE id = ?;", id)

	var strct StructInfo
	if err := row.Scan(
		&strct.Id,
		&strct.Name,
		&strct.Description,
		&strct.District,
		&strct.Region,
		&strct.Address,
		&strct.Type,
		&strct.State,
		&strct.Area,
		&strct.Owner,
		&strct.Actual_user,
		&strct.Gid,
		&strct.Permissions,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoStruct
		} else {
			return nil, err
		}
	}

	return &strct, nil
}
