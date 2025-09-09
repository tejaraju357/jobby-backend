package models

type UpdateCandidateUser struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password 	 string `json:"password"`
	ResumeUrl    string `json:"resumeurl"`
	Skills       string `json:"skills"`
}