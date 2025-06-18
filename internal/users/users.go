package users

type User struct {
	ID       int      `db:"ID"`
	Login    string   `db:"login"  json:"login"`
	Name     string   `db:"name"`
	Password Password `json:"password"`
	Hash     string   `db:"password"`
	Token    string   `db:"token"`
}

func NewUserLogin(login string, password Password) User {
	return User{Login: login, Password: password}
}

func (u *User) Authenticate() (bool, error) {
	u.Password.Verify()
}

type ConflictUserError struct {
	ID    int
	login string
}

func (e *ConflictUserError) Error() string {
	return "User with login " + e.login + " already exsists"
}

func NewConflictUserError(ID int, login string) *ConflictUserError {
	return &ConflictUserError{ID, login}
}
