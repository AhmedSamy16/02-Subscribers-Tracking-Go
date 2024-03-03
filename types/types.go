package types

import "github.com/google/uuid"

type Subscriber struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Channels []string  `json:"channels"`
}

type CreateSubscriber struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
}

type UpdateSubscriber struct {
	Name string `json:"name"`
}

type AddChannelParameters struct {
	Channel string `json:"channel"`
}

type CreatedUserResponse struct {
	UserId uuid.UUID `json:"userId"`
}
