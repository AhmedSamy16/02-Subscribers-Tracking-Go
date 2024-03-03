package application

import (
	"github.com/AhmedSamy16/02-Subscribers-API/handlers"
	"github.com/AhmedSamy16/02-Subscribers-API/repository"
	"github.com/go-chi/chi/v5"
)

func (app *App) LoadRoutes() {
	router := chi.NewRouter()

	router.Route("/subscribers", app.LoadSubscriberRoutes)

	app.router = router
}

func (app *App) LoadSubscriberRoutes(router chi.Router) {
	repo := &repository.SubscriberRepository{
		DB: app.DB,
	}
	subscriberHandler := &handlers.SubscribersHandler{
		Repo: repo,
	}

	router.Get("/", subscriberHandler.GetAllSubscribers)
	router.Post("/", subscriberHandler.CreateSubscriber)
	router.Get("/{id}", subscriberHandler.GetSubscriberById)
	router.Put("/{id}", subscriberHandler.UpdateSubscriber)
	router.Delete("/{id}", subscriberHandler.DeleteSubscriber)
	router.Post("/{id}/channel", subscriberHandler.AddChannelToUser)
}
