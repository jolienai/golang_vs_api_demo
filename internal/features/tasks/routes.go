package tasks

import "github.com/go-chi/chi/v5"

func RegisterRoutes(router chi.Router, store Store) {
	router.Route("/tasks", func(r chi.Router) {
		r.Post("/", handleCreateTask(store))
		r.Get("/", handleListTasks(store))
		r.Get("/{id}", handleGetTask(store))
		r.Patch("/{id}", handleUpdateTask(store))
		r.Delete("/{id}", handleDeleteTask(store))
	})
}
