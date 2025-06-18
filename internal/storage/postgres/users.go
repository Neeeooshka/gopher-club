package postgres

import (
	"github.com/Neeeooshka/gopher-club/internal/services/users"
)

func (l *Postgres) GetUserByLogin(login string) (users.User, error) {

	var user users.User

	row := l.DB.QueryRow("select * from gopher_users where login = $1", login)
	err := row.Scan(&user)
	if err != nil {
		return user, err
	}

	row = l.DB.QueryRow("select * from gopher_user_options where user_id = $1", user.ID)

	return user, row.Scan(&user)
}

func (l *Postgres) AddUser(user users.User, salt string) error {

	var id int
	var isNew bool

	row := l.DB.QueryRow("WITH ins AS (\n    INSERT INTO gopher_users (login, password)\n    VALUES ($1, $2)\n    ON CONFLICT (login) DO NOTHING\n        RETURNING id\n)\nSELECT id, 1 as is_new FROM ins\nUNION  ALL\nSELECT id, 0 as is_new FROM gopher_users WHERE login = $1\nLIMIT 1", user.Login, user.Password)
	err := row.Scan(&id, &isNew)
	if err != nil {
		return err
	}

	if !isNew {
		return users.NewConflictUserError(id, user.Login)
	}

	_, err = l.DB.Exec("INSERT INTO gopher_user_options (user_id, credentials) VALUES ($1, $2)", id, salt)

	return err
}
