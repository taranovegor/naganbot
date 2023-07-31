package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type User struct {
	ID        int64 `gorm:"primary_key;auto_increment:false"`
	FirstName string
	LastName  string
	Username  string
}

type UserRepository interface {
	Exists(int64) bool
	Get(int64) (User, error)
	Store(*User) error
	Update(*User) error
	GetTopShotUsersInChat(chatID int64) ([]*User, error)
	GetTopShotUsersInChatByYear(chatID int64, year int) ([]*User, error)
}

func NewUser(id int64, firstName string, lastName string, username string) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
	}
}

func (usr User) Name() string {
	var name string
	if len(usr.Username) > 0 {
		name = usr.Username
	} else {
		name = strings.Trim(fmt.Sprintf("%s %s", usr.FirstName, usr.LastName), " ")
	}

	if len(name) == 0 {
		name = strconv.FormatInt(usr.ID, 10)
	}

	return name
}

func (usr User) Mention() string {
	if len(usr.Username) > 0 {
		return fmt.Sprintf("@%s", usr.Username)
	}

	return fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", usr.ID, usr.Name())
}
