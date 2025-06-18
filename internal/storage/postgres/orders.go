package postgres

import (
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
)

func (l *Postgres) AddOrder(number string, userId int) error {

	var id int
	var isNew bool
	var uID int

	row := l.DB.QueryRow("WITH ins AS (\n    INSERT INTO gopher_orders (user_id, number)\n    VALUES ($1, $2)\n    ON CONFLICT (number) DO NOTHING\n        RETURNING user_id\n)\nSELECT id, 1 as is_new, $1 as user_id FROM ins\nUNION  ALL\nSELECT id, 0 as is_new, user_id FROM gopher_users WHERE number = $2\nLIMIT 1", userId, number)
	err := row.Scan(&id, &isNew, &uID)
	if err != nil {
		return err
	}

	if !isNew {
		if uID == userId {
			return orders.NewConflictOrderError(number)
		}
		return orders.NewConflictOrderUserError(uID, number)
	}

	return nil
}

func (l *Postgres) AddOUser(user users.User, salt string) error {

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
