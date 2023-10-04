package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/r4chi7/aspire-lite/controller"
	"github.com/r4chi7/aspire-lite/database"
	"github.com/r4chi7/aspire-lite/repository"
	"github.com/r4chi7/aspire-lite/service"
	"github.com/r4chi7/aspire-lite/token"
)

func Init() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	db := database.Get()
	userRepository := repository.NewUser(db)
	loanRepository := repository.NewLoan(db)
	loanRepaymentRepository := repository.NewLoanRepayment(db)

	userController := controller.NewUser(service.NewUser(userRepository))
	loanController := controller.NewLoan(service.NewLoan(loanRepository, loanRepaymentRepository))
	adminController := controller.NewAdmin(service.NewLoan(loanRepository, loanRepaymentRepository))

	r.Mount("/swagger", httpSwagger.WrapHandler)
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userController.Create)
		r.Post("/login", userController.Login)

		r.Route("/loans", func(r chi.Router) {
			r.Use(jwtauth.Verifier(token.GetTokenAuth()))
			r.Use(jwtauth.Authenticator)

			r.Post("/", loanController.Create)
			r.Get("/", loanController.Get)

			r.Post("/{loanID}/repay", loanController.Repay)
		})
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(jwtauth.Verifier(token.GetTokenAuth()))
		r.Use(jwtauth.Authenticator)
		r.Patch("/loan/status", adminController.UpdateLoanStatus)
	})

	return r
}
