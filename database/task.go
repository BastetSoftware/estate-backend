package database

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
)

var ErrTaskExists = errors.New("task already exists")

type Task struct {
	Id          int64
	Name        string
	Description string
	Deadline    int64
	Status      string

	Object int64

	Maintainer  int64
	Gid         int64
	Permissions uint8
}

func CreateTask(db *sql.DB, task *Task) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO tasks(name,description,deadline,status,object,maintainer,gid,permissions) VALUES(?,?,?,?,?,?,?,?);",
		task.Name,
		task.Description,
		task.Deadline,
		task.Status,
		task.Object,
		task.Maintainer,
		task.Gid,
		task.Permissions,
	)
	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == 1062 {
				return 0, ErrTaskExists
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
