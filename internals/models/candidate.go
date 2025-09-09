package models

import (
	"time"

	"gorm.io/gorm"
)

type Candidate struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	ResumeUrl    string `json:"resumeurl"`
	Skills       string `json:"skills"`
	Role         string `json:"role"`
	OTP          string
	OTPExpiry    time.Time
	Applications []Application `gorm:"foreignKey:CandidateID;references:ID"`
}
