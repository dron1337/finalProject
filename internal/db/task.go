package db

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/dron1337/finalProject/internal/models"
)

func AddTask(task *models.Task) (string, error) {
	var id int64
	res, err := DB.Exec("INSERT INTO scheduler (comment,repeat,title,date) VALUES (:comment,:repeat,:title,:date)",
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("title", task.Title),
		sql.Named("date", task.Date))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return strconv.Itoa(int(id)), err
}

func GetTasks(limit int) (*models.TasksResp, error) {
	tasks := make([]models.Task, 0)
	rows, err := DB.Query("SELECT id,date,title,comment,repeat FROM scheduler ORDER BY date ASC LIMIT :limit",
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &models.TasksResp{Tasks: tasks}, nil
}

func GetTask(id string) (*models.Task, error) {
	t := &models.Task{}
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func UpdateTask(task *models.Task) error {
	res, err := DB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
		sql.Named("date", next),
		sql.Named("id", id))
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("error delete record")
	}
	return nil
}
