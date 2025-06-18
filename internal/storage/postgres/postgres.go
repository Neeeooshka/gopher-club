package postgres

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

type ConflictError struct {
	ShortLink string
}

func (e *ConflictError) Error() string {
	return "link already exsists"
}

type Postgres struct {
	DB *sql.DB
}

func (l *Postgres) Close() error {
	return l.DB.Close()
}

func (l *Postgres) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := l.DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (l *Postgres) initStructForLinks() (err error) {
	_, err = l.DB.Exec("CREATE TABLE IF NOT EXISTS shortener_links (\n    id SERIAL,\n    short_url character(8) NOT NULL,\n    original_url character varying(250) NOT NULL,\n    user_id character(32) NULL,\n    deleted boolean NOT NULL DEFAULT false,\n    PRIMARY KEY (id),\n    UNIQUE (original_url)\n )")
	return err
}

func NewPostgresLinksStorage(conn string) (pgx *Postgres, err error) {

	pgx = &Postgres{}

	pgx.DB, err = sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}

	return pgx, pgx.initStructForLinks()
}
