package models

import "time"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User_id struct {
	Id string `json:"id"`
}

type Session struct {
	User_id string    `json:"user_id"`
	Expiry  time.Time `json:"expiry"`
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
