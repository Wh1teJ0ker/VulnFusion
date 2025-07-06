package storage

type User struct {
	ID       int
	Username string
	Password string
}

func UserExists(username string) bool {
	row := DB.QueryRow("SELECT 1 FROM users WHERE username = ?", username)
	var tmp int
	return row.Scan(&tmp) == nil
}

func CreateUser(username, password string) error {
	_, err := DB.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, password)
	return err
}

func GetUserByUsername(username string) (*User, error) {
	row := DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Password)
	return &u, err
}
