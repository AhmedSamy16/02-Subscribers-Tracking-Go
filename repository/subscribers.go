package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/AhmedSamy16/02-Subscribers-API/types"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type SubscriberRepository struct {
	DB *sql.DB
}

func (repo *SubscriberRepository) GetAllSubscribers(ctx context.Context) (*[]types.Subscriber, error) {
	subStmt, err := repo.DB.PrepareContext(ctx, `
		SELECT s.id, s.name, json_agg(c.title) AS channels
		FROM subscribers AS s
		LEFT JOIN subscriber_channels AS c 
		ON s.id = c.user_id
		GROUP BY s.id
	`)
	if err != nil {
		return nil, errors.New("couldn't prepare statement")
	}
	defer subStmt.Close()

	rows, err := subStmt.QueryContext(ctx)
	if err != nil {
		return nil, errors.New("couldn't execute statement")
	}
	defer rows.Close()

	var userId uuid.UUID
	var name string
	var channelsJson string

	subs := []types.Subscriber{}

	for rows.Next() {
		err := rows.Scan(&userId, &name, &channelsJson)
		if err != nil {
			return nil, errors.New("couldn't read data")
		}
		var channels []string
		if err = json.Unmarshal([]byte(channelsJson), &channels); err != nil {
			return nil, errors.New("couldn't read channels")
		}
		sub := types.Subscriber{
			Id:       userId,
			Name:     name,
			Channels: channels,
		}

		subs = append(subs, sub)
	}

	return &subs, nil
}

func (repo *SubscriberRepository) GetSubscriberById(ctx context.Context, id uuid.UUID) (*types.Subscriber, error) {
	subStmt, err := repo.DB.PrepareContext(ctx, `
		SELECT s.id, s.name, json_agg(c.title) AS channels
		FROM subscribers AS s
		LEFT JOIN subscriber_channels AS c 
		ON s.id = c.user_id
		WHERE s.id = $1
		GROUP BY s.id
	`)
	if err != nil {
		log.Println(err)
		return nil, errors.New("couldn't prepare statement")
	}
	defer subStmt.Close()

	rows, err := subStmt.QueryContext(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("couldn't execute statement")
	}
	defer rows.Close()

	var userId uuid.UUID
	var name string
	var channelsJson string

	var channels []string

	if rows.Next() {
		if err = rows.Scan(&userId, &name, &channelsJson); err != nil {
			log.Println(err)
			return nil, errors.New("couldn't read data")
		}
		if err = json.Unmarshal([]byte(channelsJson), &channels); err != nil {
			log.Println(err)
			return nil, errors.New("couldn't read channels")
		}
	} else {
		return nil, nil
	}

	return &types.Subscriber{
		Id:       userId,
		Name:     name,
		Channels: channels,
	}, nil
}

func (repo *SubscriberRepository) CreateSubscriber(ctx context.Context, data types.CreateSubscriber) (*uuid.UUID, error) {
	// 1. Start Transaction
	tx, err := repo.DB.Begin()
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to start transaction")
	}

	// 2. Insert User
	userStmt, err := tx.Prepare("INSERT INTO subscribers (id, name) VALUES ($1, $2) RETURNING id")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer userStmt.Close()

	var userId uuid.UUID
	err = userStmt.QueryRowContext(ctx, uuid.New(), data.Name).Scan(&userId)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to insert user")
	}

	// 3. Insert Channels
	channelsStmt, err := tx.Prepare("INSERT INTO subscriber_channels (user_id, title) VALUES ($1, $2)")
	if err != nil {
		tx.Rollback()
		return nil, errors.New("failed to prepare channels")
	}
	defer channelsStmt.Close()
	for _, channel := range data.Channels {
		_, err := channelsStmt.Exec(userId, channel)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("failed to insert channel")
		}
	}

	// 4. Done
	err = tx.Commit()
	if err != nil {
		return nil, errors.New("failed to commit")
	}
	return &userId, nil
}

func (repo *SubscriberRepository) UpdateSubscriber(ctx context.Context, id uuid.UUID, data *types.UpdateSubscriber) error {
	stmt, err := repo.DB.PrepareContext(ctx, "UPDATE subscribers SET name = $1 WHERE id = $2")
	if err != nil {
		return errors.New("failed to create statement")
	}
	_, err = stmt.ExecContext(ctx, data.Name, id)
	if err != nil {
		return errors.New("failed to update subscriber")
	}
	return nil
}

func (repo *SubscriberRepository) DeleteSubscriber(ctx context.Context, id uuid.UUID) error {
	stmt, err := repo.DB.PrepareContext(ctx, "DELETE FROM subscribers WHERE id = $1")
	if err != nil {
		return errors.New("failed to create statement")
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return errors.New("failed to delete subscriber")
	}
	return nil
}

func (repo *SubscriberRepository) AddChannelToSubscriber(ctx context.Context, id uuid.UUID, data *types.AddChannelParameters) error {
	if data.Channel == "" {
		return errors.New("invalid data")
	}
	stmt, err := repo.DB.PrepareContext(ctx, "INSERT INTO subscriber_channels (user_id, title) VALUES ($1, $2)")
	if err != nil {
		return errors.New("failed to prepare statement")
	}
	_, err = stmt.Exec(id, data.Channel)
	if err != nil {
		return errors.New("failed to insert data")
	}
	return nil
}
