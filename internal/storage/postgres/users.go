package postgres

import "github.com/Neeeooshka/gopher-club.git/internal/auth"

type ConflictUserError struct {
	ID    int
	login string
}

func (e *ConflictUserError) Error() string {
	return "User with login " + e.login + " already exsists"
}

func (l *Postgres) GetUserByLogin(login string) (*auth.User, bool) {

	user := auth.User{}

	row := l.DB.QueryRow("select * from gopher_users where login = $1", login)
	err := row.Scan(&user)
	if err != nil {
		return &auth.User{}, false
	}

	return &user, true
}

func (l *Postgres) AddUser(user auth.User) error {

	var id int
	var isNew bool

	row := l.DB.QueryRow("WITH ins AS (\n    INSERT INTO gopher_users (login, password, name)\n    VALUES ($1, $2, $3)\n    ON CONFLICT (login) DO NOTHING\n        RETURNING id\n)\nSELECT id, 1 as is_new FROM ins\nUNION  ALL\nSELECT id, 0 as is_new FROM gopher_users WHERE login = $1\nLIMIT 1", user.Login, user.Password.GetHash(), user.Name)
	err := row.Scan(&id, &isNew)
	if err != nil {
		return err
	}

	if !isNew {
		return &ConflictUserError{ID: id, login: user.Login}
	}

	return nil
}
