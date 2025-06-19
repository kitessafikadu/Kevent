package database

import (
	"database/sql"
	"context"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	Id          int  `json:"id"`
	OwnerId    int  `json:"ownerId" binding:"required"`
	Name 			string `json:"name" binding:"required, min=3"` 
	Description string `json:"description" binding:"required, min=10"`
	Location    string `json:"location" binding:"required, min=3"`
	Date				string `json:"date" binding:"required, datetime=2006-01-02"`
}

func (m *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO events (owner_id, name, description, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location).Scan(&event.Id)
}

func (m *EventModel) GetAll()([]*Event, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM events"
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (m *EventModel) Get(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM events WHERE id = $1"

	var event Event

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}