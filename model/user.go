package model

import (
	"time"
)

type User struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	DateTime time.Time `json:"date_time"`
}

func (m *User) TableName() string {
	return "user"
}
