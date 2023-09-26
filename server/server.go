package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/r4chi7/aspire-lite/controller"
	"github.com/r4chi7/aspire-lite/database"
	"github.com/r4chi7/aspire-lite/repository"
	"github.com/r4chi7/aspire-lite/service"
)

func Init() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Logger)

	db := database.Get()
	userRepository := repository.NewUser(db)

	userController := controller.NewUser(service.NewUser(userRepository))

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userController.Create)
		r.Post("/login", userController.Login)
	})

	return r
}
