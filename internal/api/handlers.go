package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dron1337/finalProject/internal/constants"
	"github.com/dron1337/finalProject/internal/db"
	"github.com/dron1337/finalProject/internal/models"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	var task models.Task
	if err := json.Unmarshal(body, &task); err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := task.Validate(); err != nil {

		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	writeJSON(w, models.TaskResponse{ID: id}, http.StatusOK)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, models.TaskResponse{Error: "id parameter is required"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	writeJSON(w, task, http.StatusOK)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasks(50)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: "Task not found"}, http.StatusNotFound)
		return
	}
	writeJSON(w, tasks, http.StatusOK)
}

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	repeat := r.FormValue("repeat")
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	var err error
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(constants.DateFormat, nowStr)
		if err != nil {
			writeJSON(w, models.TaskResponse{Error: "Invalid 'now' parameter format, expected YYYY-MM-DD"}, http.StatusBadRequest)
			return
		}
	}

	d, err := models.NextDate(now, dateStr, repeat)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(d))
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.Unmarshal(body, &task); err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	if err := task.Validate(); err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	writeJSON(w, task, http.StatusOK)
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, models.TaskResponse{Error: "id parameter is required"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	if len(task.Repeat) == 0 {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
			return
		}
	} else {
		d, err := models.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		err = db.UpdateDate(d, id)
		if err != nil {
			writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusInternalServerError)
			return
		}
	}
	writeJSON(w, models.TaskResponse{}, http.StatusOK)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJSON(w, models.TaskResponse{}, http.StatusOK)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendAuthError(w, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	var auth models.Authorization
	if err := json.Unmarshal(body, &auth); err != nil {
		writeJSON(w, models.TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	secret := os.Getenv("TODO_PASSWORD")
	if secret == "" {
		writeJSON(w, models.TaskResponse{Error: "Authentication not configured"}, http.StatusInternalServerError)
		return
	}

	if auth.Password != secret {
		writeJSON(w, models.TaskResponse{Error: "Incorrect password"}, http.StatusUnauthorized)
		return
	}

	token, err := generateToken(secret)
	if err != nil {
		writeJSON(w, models.TaskResponse{Error: "Failed to generate token"}, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(tokenExpiration.Seconds()),
		SameSite: http.SameSiteStrictMode,
	})

	writeJSON(w, models.TaskResponse{}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
