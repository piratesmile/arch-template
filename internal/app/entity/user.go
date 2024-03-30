package entity

import "arch-template/ent"

type User struct {
	ID       uint
	UserName string
	Password string
}

func UserFromEntModel(u *ent.User) *User {
	return &User{
		ID:       u.ID,
		UserName: u.Username,
		Password: u.Password,
	}
}
