package models

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Salary      string    `json:"salary"`
	CompanyID   uint      `json:"companyid"`
	Company   Company   `gorm:"foreignKey:CompanyID"`
	Applications []Application `gorm:"foreignKey:JobID"`
}