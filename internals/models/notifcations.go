package models 

import "gorm.io/gorm" 

type Notifications struct {
	gorm.Model
	UserID    uint      `json:"user_id"`     // candidate or company user
	UserType  string    `json:"user_type"`   // "candidate" or "company"
	Message   string    `json:"message"`
}