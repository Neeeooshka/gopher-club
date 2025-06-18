package postgres

import "github.com/Neeeooshka/gopher-club/internal/users"

func (l *Postgres) GetUserByLogin(login string) (users.User, error) {

	var user users.User

	row := l.DB.QueryRow("select * from gopher_users where login = $1", login)

	return user, row.Scan(&user)
}

func (l *Postgres) GetUserKey(ID int) (string, error) {

	var key string

	row := l.DB.QueryRow("select * from gopher_keys where ID = $1", ID)

	return key, row.Scan(&key)
}

func (l *Postgres) AddUser(user users.User) error {

	var id int
	var isNew bool

	row := l.DB.QueryRow("WITH ins AS (\n    INSERT INTO gopher_users (login, password, name)\n    VALUES ($1, $2, $3)\n    ON CONFLICT (login) DO NOTHING\n        RETURNING id\n)\nSELECT id, 1 as is_new FROM ins\nUNION  ALL\nSELECT id, 0 as is_new FROM gopher_users WHERE login = $1\nLIMIT 1", user.Login, user.Password.GetHash(), user.Name)
	err := row.Scan(&id, &isNew)
	if err != nil {
		return err
	}

	if !isNew {
		return users.NewConflictUserError(id, user.Login)
	}

	return nil
}
