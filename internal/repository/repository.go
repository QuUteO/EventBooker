package repository

import pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"

type EventRepository struct {
	conn *pgxdriver.Postgres
}

func New(conn *pgxdriver.Postgres) *EventRepository {
	return &EventRepository{conn: conn}
}
