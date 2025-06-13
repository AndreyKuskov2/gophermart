package storage

const (
	createNewUser          = "INSERT INTO users (login, password) VALUES ($1, $2);"
	checkUserIsExists      = "SELECT user_id FROM users WHERE login = $1;"
	getUserPasswordByLogin = "SELECT password FROM users WHERE login = $1;"
)
