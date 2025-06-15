package storage

const (
	// register and login
	createNewUser          = "INSERT INTO users(login, password) VALUES ($1, $2) RETURNING user_id;"
	checkUserIsExists      = "SELECT user_id FROM users WHERE login = $1;"
	getUserPasswordByLogin = "SELECT user_id, password FROM users WHERE login = $1;"

	//
	createOrder       = "INSERT INTO orders(number, status, accrual, user_id) VALUES ($1, $2, $3, $4);"
	getOrderByNumber  = "SELECT * FROM orders WHERE number = $1;"
	getOrdersByUserID = "SELECT * FROM orders WHERE user_id = $1;"
	getUserBalance    = `SELECT
	  COALESCE(accrual_sum, 0) - COALESCE(withdrawn_sum, 0) AS current,
	  COALESCE(withdrawn_sum, 0) AS withdrawn
	FROM
	  (SELECT SUM(accrual) AS accrual_sum FROM orders WHERE user_id = $1 AND status = $2) o,
	  (SELECT SUM(amount) AS withdrawn_sum FROM withdrawals WHERE user_id = $1) w`
)
