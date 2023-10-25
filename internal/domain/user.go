package domain

import "time"

// User 领域对象
type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}

func (u User) Encrypt() {

}

func (u User) ComparePassword() {

}
