package models

type User struct {
	id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
