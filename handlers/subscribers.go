package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AhmedSamy16/02-Subscribers-API/repository"
	"github.com/AhmedSamy16/02-Subscribers-API/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubscribersHandler struct {
	Repo *repository.SubscriberRepository
}

func (handler *SubscribersHandler) GetAllSubscribers(w http.ResponseWriter, r *http.Request) {
	data, err := handler.Repo.GetAllSubscribers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}
	respondWithJSON(w, 200, data)
}

func (handler *SubscribersHandler) GetSubscriberById(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	subId, err := uuid.Parse(params)
	if err != nil {
		respondWithError(w, 400, "Invalid id")
		return
	}
	data, err := handler.Repo.GetSubscriberById(r.Context(), subId)
	if err != nil {
		respondWithError(w, 500, "Something went wrong, Please try again")
		return
	}
	if data == nil {
		respondWithError(w, 404, "Subscriber not found")
		return
	}
	respondWithJSON(w, 200, data)
}

func (handler *SubscribersHandler) CreateSubscriber(w http.ResponseWriter, r *http.Request) {
	sub := types.CreateSubscriber{}
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		respondWithError(w, 400, "Invalid input")
		return
	}
	defer r.Body.Close()

	res, err := handler.Repo.CreateSubscriber(r.Context(), sub)
	if err != nil {
		respondWithError(w, 400, err.Error())
	} else {
		respondWithJSON(w, 200, types.CreatedUserResponse{
			UserId: *res,
		})
	}
}

func (handler *SubscribersHandler) AddChannelToUser(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	subId, err := uuid.Parse(params)
	if err != nil {
		respondWithError(w, 400, "Invalid id")
		return
	}
	var channel types.AddChannelParameters
	if err = json.NewDecoder(r.Body).Decode(&channel); err != nil {
		respondWithError(w, 400, "Invalid input")
		return
	}
	err = handler.Repo.AddChannelToSubscriber(r.Context(), subId, &channel)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	respondWithJSON(w, 201, struct{}{})

}

func (handler *SubscribersHandler) UpdateSubscriber(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	subId, err := uuid.Parse(params)
	if err != nil {
		respondWithError(w, 400, "Invalid id")
		return
	}
	var data types.UpdateSubscriber
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, 400, "Invalid input")
		return
	}

	err = handler.Repo.UpdateSubscriber(r.Context(), subId, &data)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 200, struct{}{})
}

func (handler *SubscribersHandler) DeleteSubscriber(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	subId, err := uuid.Parse(params)
	if err != nil {
		respondWithError(w, 400, "Invalid id")
		return
	}
	err = handler.Repo.DeleteSubscriber(r.Context(), subId)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	respondWithJSON(w, 204, struct{}{})
}
