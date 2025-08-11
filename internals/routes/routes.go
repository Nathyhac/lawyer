package routes

import (
	app "github.com/Nathac/go-api/internals"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/lawyer", app.LawyerHandler.GetAllLawyers)
	r.Get("/lawyer/{id}", app.LawyerHandler.GetLawyerById)
	r.Post("/lawyer", app.LawyerHandler.CreateLawyer)
	r.Put("/lawyer", app.LawyerHandler.UpdateLawyer)
	r.Delete("/lawyer/{id}", app.LawyerHandler.Deletelawyer)
	r.Post("/user", app.UserHandler.HandlerUserRegister)
	r.Get("/user/{username}", app.UserHandler.GetUserByUsername)
	r.Post("/user/login", app.TokenHandler.CreateTokenHandler)
	return r
}
