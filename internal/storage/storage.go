package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ega-forever/otus-image-service/internal/domain/models"
	"github.com/jackc/pgx/v4"
	"log"
)

type Storage struct {
	ctx context.Context
	db  *pgx.Conn
}

func New(uri string) *Storage {

	ctx := context.Background()
	db, err := pgx.Connect(ctx, uri) // *sql.DB
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	return &Storage{ctx: ctx, db: db}
}

func (storage *Storage) Migrate() error {

	query := `create table if not exists events (
    id SERIAL,
    text text,
    timestamp double precision
)`

	_, err := storage.db.Exec(storage.ctx, query) // sql.Result
	if err != nil {
		return err
	}

	fmt.Println("migrated tables")
	return nil
}

func (storage *Storage) SaveEvent(ctx context.Context, event *models.Event) (*models.Event, error) {

	query := `insert into events(text, timestamp) values($1, $2) RETURNING id`
	// statement, _ := storage.db.Prepare(storage.ctx, "save_event", query)

	var id int64
	err := storage.db.QueryRow(storage.ctx, query, event.Text, event.Timestamp).Scan(&id)

	if err != nil {
		return nil, err
	}

	ev := models.Event{Id: id, Text: event.Text, Timestamp: event.Timestamp}

	return &ev, nil
}

func (storage *Storage) UpdateEventById(ctx context.Context, id int64, text string, timestamp int64) (*models.Event, error) {

	query := `update events set text = $1, timestamp = $2 where id = $3`
	result, err := storage.db.Exec(storage.ctx, query, text, timestamp, id)

	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, errors.New("record not found")
	}

	return &models.Event{Id: id, Text: text, Timestamp: timestamp}, nil
}

func (storage *Storage) GetEventById(ctx context.Context, id int64) (*models.Event, error) {

	query := "select * from events where id = $1"
	row := storage.db.QueryRow(storage.ctx, query, id)

	ev := models.Event{}

	err := row.Scan(&ev.Id, &ev.Text, &ev.Timestamp)

	if err == sql.ErrNoRows {
		return nil, errors.New("record not found")
	} else if err != nil {
		return nil, err
	}

	return &ev, nil
}

func (storage *Storage) DeleteEventById(ctx context.Context, id int64) error {

	query := `delete from events where id = $1`
	result, err := storage.db.Exec(storage.ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (storage *Storage) ListEvents(ctx context.Context) ([]*models.Event, error) {

	query := `select * from events`
	rows, err := storage.db.Query(storage.ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	envs := make([]*models.Event, 0)

	for rows.Next() {
		ev := models.Event{}
		err := rows.Scan(&ev.Id, &ev.Text, &ev.Timestamp)

		if err != nil {
			return nil, err
		}

		envs = append(envs, &ev)
	}

	return envs, nil
}
