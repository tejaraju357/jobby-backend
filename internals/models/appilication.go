package models

import (
	"gorm.io/gorm"
)


type Application struct {
    gorm.Model
    CandidateID uint   `json:"candidateid"`
    JobID       uint   `json:"jobid"`
    Status      string `json:"status"`

    Candidate   Candidate `gorm:"foreignKey:CandidateID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    Job         Job       `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
