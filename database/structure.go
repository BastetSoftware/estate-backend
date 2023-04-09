package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

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

type ArgsFStructFind struct {
	Token       string
	Name        string
	Description string
	District    string
	Region      string
	Address     string
	Type        string
	State       string
	AreaFrom    *int32
	AreaTo      *int32
	Owner       string
	Actual_user string
	Gid         *int64
	Limit       int16
	SortAsc     bool
	Offset      int16
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

	strct.Id, err = result.LastInsertId()
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

func FindStructures(db *sql.DB, filter ArgsFStructFind) ([]StructInfo, error) {
	query := "SELECT * FROM objects WHERE "
	var params []string
	if filter.Name != "" {
		params = append(params, "name = \""+filter.Name+"\"")
	}
	if filter.Description != "" {
		params = append(params, "desciprtion = \""+filter.Description+"\"")
	}
	if filter.District != "" {
		params = append(params, "district = \""+filter.District+"\"")
	}
	if filter.Region != "" {
		params = append(params, "region = \""+filter.Region+"\"")
	}
	if filter.Address != "" {
		params = append(params, "address = \""+filter.Address+"\"")
	}
	if filter.Type != "" {
		params = append(params, "type = \""+filter.Type+"\"")
	}
	if filter.State != "" {
		params = append(params, "state= \""+filter.State+"\"")
	}
	if filter.AreaFrom != nil && filter.AreaTo != nil {
		params = append(params, "area >= "+strconv.FormatInt(int64(*filter.AreaFrom), 10)+" and "+"area <= "+strconv.FormatInt(int64(*filter.AreaTo), 10))
	} else if filter.AreaTo != nil {
		params = append(params, "area <"+strconv.FormatInt(int64(*filter.AreaTo), 10))
	} else if filter.AreaFrom != nil {
		params = append(params, "area > "+strconv.FormatInt(int64(*filter.AreaFrom), 10))
	}
	if filter.Gid != nil {
		params = append(params, "gid = "+strconv.FormatInt(int64(*filter.Gid), 10))
	}
	for i := 0; i < len(params); i++ {
		if i == 0 {
			query += params[i]
		} else {
			query += " AND " + params[i]
		}
	}
	fmt.Println(query)
	rows, err := db.Query(query + " LIMIT " + strconv.FormatInt(int64(filter.Limit), 10))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	structures := make([]StructInfo, 0)
	for rows.Next() {
		t := StructInfo{}
		err := rows.Scan(
			&t.Id,
			&t.Name,
			&t.Description,
			&t.District,
			&t.Region,
			&t.Address,
			&t.Type,
			&t.State,
			&t.Area,
			&t.Owner,
			&t.Actual_user,
			&t.Gid,
			&t.Permissions,
		)
		if err != nil {
			return nil, err
		}
		structures = append(structures, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return structures, nil

}

func DeleteStruct(db *sql.DB, Id int64) error {
	_, err := db.Exec(
		"DELETE FROM objects WHERE id=?;",
		Id,
	)
	if err != nil {
		return err
	}
	return nil
}
