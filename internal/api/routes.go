package api

import (
	"net/http"

	"github.com/dron1337/finalProject/internal/constants"
	"github.com/go-chi/chi"
)

func Init() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/nextdate", NextDayHandler)
	r.Post("/api/signin", SignInHandler)
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Route("/api/task", func(r chi.Router) {
			r.Post("/", AddTaskHandler)
			r.Post("/done", DoneTaskHandler)
			r.Get("/", GetTaskHandler)
			r.Delete("/", DeleteTaskHandler)
			r.Put("/", UpdateTaskHandler)
		})

		r.Get("/api/tasks", GetTasksHandler)
	})

	if constants.WebDir != "" {
		fs := http.FileServer(http.Dir(constants.WebDir))
		r.Handle("/*", fs)
	}

	return r
}
