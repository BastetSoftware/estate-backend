package database

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
)

var ErrTaskExists = errors.New("task already exists")
var ErrNoTask = errors.New("task does not exist")

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

type TaskFilter struct {
	Name         *string
	Description  *string
	DeadlineFrom *int64
	DeadlineTo   *int64
	Status       *string
	Object       *int64
	Maintainer   *int64
	Gid          *int64

	Limit  int16
	Offset int16
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

func RemoveTask(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM tasks WHERE id=?;", id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n < 1 {
		return ErrNoTask
	}

	return err
}

func GetTask(db *sql.DB, id int64) (*Task, error) {
	row := db.QueryRow("SELECT * FROM tasks WHERE id=?;", id)

	var task Task
	err := row.Scan(
		&task.Id,
		&task.Name,
		&task.Description,
		&task.Deadline,
		&task.Status,
		&task.Object,
		&task.Maintainer,
		&task.Gid,
		&task.Permissions,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoTask
		} else {
			return nil, err
		}
	}

	return &task, nil
}

func FilterTasks(db *sql.DB, filter *TaskFilter) ([]*Task, error) {
	rows, err := db.Query(`SELECT * FROM tasks
	    WHERE ((name LIKE ?) OR ? IS NULL)
	      AND ((description LIKE ?) OR ? IS NULL)
	      AND ((deadline >= ?) OR ? IS NULL)
	      AND ((deadline <= ?) OR ? IS NULL)
	      AND ((status LIKE ?) OR ? IS NULL)
	      AND ((object = ?) OR ? IS NULL)
	      AND ((maintainer = ?) OR ? IS NULL)
	      AND ((Gid = ?) OR ? IS NULL)
	    LIMIT ? OFFSET ?;`,

		filter.Name,
		filter.Name,
		filter.Description,
		filter.Description,
		filter.DeadlineFrom,
		filter.DeadlineFrom,
		filter.DeadlineTo,
		filter.DeadlineTo,
		filter.Status,
		filter.Status,
		filter.Object,
		filter.Object,
		filter.Maintainer,
		filter.Maintainer,
		filter.Gid,
		filter.Gid,

		filter.Limit,
		filter.Offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.Id,
			&task.Name,
			&task.Description,
			&task.Deadline,
			&task.Status,
			&task.Object,
			&task.Maintainer,
			&task.Gid,
			&task.Permissions,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}
