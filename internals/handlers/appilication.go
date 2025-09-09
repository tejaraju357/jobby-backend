package handlers

import (
	"jobby/internals/db"
	"jobby/internals/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CandidateSafe struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	ResumeUrl string `json:"resumeurl"`
	Skills    string `json:"skills"`
	Role      string `json:"role"`
}

type ApplicationResponse struct {
	ID        uint         `json:"id"`
	Status    string       `json:"status"`
	Candidate CandidateSafe `json:"candidate"`
	JobID     uint         `json:"job_id"`
	CreatedAt time.Time    `json:"created_at"`
}

func GetApplicationByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	var application models.Application
	if err := db.DB.Preload("Candidate").Where("id = ?", id).First(&application).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Application not found"})
	}

	candidate := CandidateSafe{
		Name:      application.Candidate.Name,
		Email:     application.Candidate.Email,
		ResumeUrl: application.Candidate.ResumeUrl,
		Skills:    application.Candidate.Skills,
		Role:      application.Candidate.Role,
	}

	resp := ApplicationResponse{
		Status:    application.Status,
		Candidate: candidate,
		JobID:     application.JobID,
		CreatedAt: application.CreatedAt,
	}

	return c.JSON(resp)
}

func ListApplicationsByCandidate(c *fiber.Ctx) error {
	candidateID := c.Params("candidate_id")
	if candidateID == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	var applications []models.Application
	if err := db.DB.Where("candidate_id = ?", candidateID).Find(&applications).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Applications not found"})
	}

	var resp []ApplicationResponse
	for _, app := range applications {
		candidate := CandidateSafe{
			ID:        app.Candidate.ID,
			Name:      app.Candidate.Name,
			Email:     app.Candidate.Email,
			ResumeUrl: app.Candidate.ResumeUrl,
			Skills:    app.Candidate.Skills,
			Role:      app.Candidate.Role,
		}

		resp = append(resp, ApplicationResponse{
			ID:        app.ID,
			Status:    app.Status,
			Candidate: candidate,
			JobID:     app.JobID,
			CreatedAt: app.CreatedAt,
		})
	}

	return c.JSON(resp)
}

func ListApplicationsByJob(c *fiber.Ctx) error {
	jobID := c.Params("job_id")
	if jobID == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	var applications []models.Application
	if err := db.DB.Preload("Candidate").Where("job_id = ?", jobID).Find(&applications).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Applications not found"})
	}
	return c.JSON(fiber.Map{"Applications": applications})
}

func UpdateApplicationStatus(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	updateStatus := &models.Application{}
	if err := c.BodyParser(updateStatus); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	result := db.DB.Model(&models.Candidate{}).Where("id = ?", id).Update("status", updateStatus.Status)

	if result.Error != nil {
		log.Printf("Unable to update database: %v\n", result.Error)
		return c.Status(500).JSON(fiber.Map{
			"error": "Cannot update database",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func DeleteApplication(c *fiber.Ctx) error {

	applicationModel := models.Application{}
	id := c.Params("id")
	if id == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	if id == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := db.DB.Delete(applicationModel, id)

	if err.Error != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "could not delete the application",
		})
		return err.Error
	}

	c.Status(200).JSON(&fiber.Map{
		"message": "application deleted successull",
	})
	return nil
}
