package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dron1337/finalProject/internal/constants"
)

type TasksResp struct {
	Tasks []Task `json:"tasks"`
}
type Task struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Date    string `json:"date"`
	Repeat  string `json:"repeat"`
}
type Authorization struct {
	Password string `json:"password"`
}

type TaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	return t.checkDate()
}

func (t *Task) checkDate() error {
	now := time.Now()
	today := now.Format(constants.DateFormat)
	if t.Date == "" {
		t.Date = now.Format(constants.DateFormat)
		return nil
	}
	parsedDate, err := time.Parse(constants.DateFormat, t.Date)
	if err != nil {
		return fmt.Errorf("invalid start date format, expected YYYYMMDD")
	}

	if t.Date == today {
		return nil
	}

	if parsedDate.Before(now) {
		if t.Repeat == "" {
			t.Date = now.Format(constants.DateFormat)
		} else {
			nextDate, err := NextDate(now, t.Date, t.Repeat)
			if err != nil {
				return err
			}
			t.Date = nextDate
		}
	}
	return nil
}

func NextDate(now time.Time, startDate string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat parameter is required")
	}

	parsedDate, err := time.Parse(constants.DateFormat, startDate)
	if err != nil {
		return "", fmt.Errorf("invalid start date format, expected YYYY-MM-DD")
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if len(parts) == 1 {
			return "", fmt.Errorf("missing day count in repeat parameter")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid day count in repeat parameter")
		}

		if days <= 0 {
			return "", fmt.Errorf("day count must be positive")
		}

		if days > 400 {
			return "", fmt.Errorf("day count exceeds maximum value of 400")
		}

		for {
			parsedDate = parsedDate.AddDate(0, 0, days)
			if parsedDate.After(now) {
				break
			}
		}

	case "y":
		for {
			parsedDate = parsedDate.AddDate(1, 0, 0)
			if parsedDate.After(now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat type, use 'd' for days or 'y' for years")
	}

	return parsedDate.Format(constants.DateFormat), nil
}
