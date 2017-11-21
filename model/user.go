package model

// User stores user information
type User struct {
	// ID of User
	ID int64 `gorm:primary key;not_nil`
	// Login of User
	Login string `gorm:not_nil`
	// Pass of User
	Pass string `gorm:not_nil`
	// WorkNumber of User
	WorkNumber int32
}

// Get fetching User info
func (u *User) Get(login, password string) error {
	DBConn.SingularTable(true)
	return DBConn.
		Where("name = ? AND pass = ?", login, password).
		First(&u).
		Error
}

// Update User information
func (u *User) Update() error {
	DBConn.SingularTable(true)
	return DBConn.
		Model(&u).
		Update("pass", &u.Pass).
		Error
}
