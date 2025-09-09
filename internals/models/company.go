package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Role        string `json:"role"`
	OTP         string
	OTPExpiry   time.Time
	Jobs        []Job `gorm:"foreignKey:CompanyID"`
}
