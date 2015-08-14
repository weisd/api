package user

import (
	"../../models"
)

type User struct {
	Id   int64
	Name string
}

func Create(u *User) (aft int64, err error) {
	x := models.Master()
	return x.Insert(u)
}

func GetName(id int64) (name string, err error) {
	x := models.Slave()

	user := &User{Id: id}
	has, err := x.Id(id).Get(user)
	if err != nil || !has {
		return
	}

	return user.Name, nil
}
