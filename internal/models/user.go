package models

type User struct {
	ID    string `json:"-"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
