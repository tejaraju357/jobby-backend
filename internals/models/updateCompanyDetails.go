package models

type UpdateCompanyDetails struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
	Website     string `json:"website"`
}