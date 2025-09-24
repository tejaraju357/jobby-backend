package handlers

import (
	"fmt"
	"jobby/internals/cache"
	"jobby/internals/db"
	"jobby/internals/models"
	"jobby/internals/utils"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterCandidate(c *fiber.Ctx) error {
	var user models.Candidate
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":   "invalid request body",
			"details": err.Error(),
		})
	}

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Password = strings.TrimSpace(user.Password)
	user.Name = strings.TrimSpace(user.Name)
	user.ResumeUrl = strings.TrimSpace(user.ResumeUrl)
	user.Skills = strings.TrimSpace(user.Skills)

	if user.Name == "" || user.Password == "" || user.Email == "" || user.ResumeUrl == "" || user.Skills == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "please fill all required details",
		})
	}

	var existing models.Candidate
	if err := db.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
			"error": "email already exists",
		})
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error":   "failed to register user",
			"details": err.Error(),
		})
	}
	PublishNotification(user.ID, "candidate", "🎉 You have successfully registered!")
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message":   "candidate created successfully",
		"candidate": user,
	})
}

func LoginCandidate(c *fiber.Ctx) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var err error
	err = c.BodyParser(&credentials)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":   "inavlid request body",
			"details": err.Error(),
		})
		return nil
	}

	if credentials.Email == "" || credentials.Password == "" {
		c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "email and password required",
		})
		return nil
	}

	var user models.Candidate

	credentials.Email = strings.TrimSpace(strings.ToLower(credentials.Email))
	credentials.Password = strings.TrimSpace(credentials.Password)

	fmt.Println("Login attempt:", credentials.Email, credentials.Password)
	err = db.DB.Where("email = ? AND password=?", credentials.Email, credentials.Password).First(&user).Error
	fmt.Println("DB query error:", err)
	fmt.Println("DB user found:", user.Email, user.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Invalid email and password",
		})
	}

	otp := fmt.Sprintf("%05d", rand.Intn(1000))
	expiry := time.Now().Add(5 * time.Minute)
	user.OTP = otp
	user.OTPExpiry = expiry
	db.DB.Save(&user)

	if err := utils.SendOtp(user.Email, otp); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to send the otp",
		})
	}

	return c.JSON(&fiber.Map{
		"message": "Otp is sented to your Email",
	})
}

func VerifyOTPCandidate(c *fiber.Ctx) error {

	var payload struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	err := c.BodyParser(&payload)
	if err != nil {
		c.Status(402).JSON(&fiber.Map{
			"error": "invalid request",
		})
		return nil
	}

	var user models.Candidate
	err = db.DB.Where("email = ?", payload.Email).First(&user).Error
	if err != nil {
		c.Status(402).JSON(&fiber.Map{
			"error": "invalid email",
		})
		return nil
	}

	if user.OTP != payload.OTP || time.Now().After(user.OTPExpiry) {
		c.Status(402).JSON(&fiber.Map{
			"error": "invalid or expired token",
		})
		return nil
	}

	user.OTP = ""
	user.OTPExpiry = time.Time{}
	db.DB.Save(&user)

	token, err := utils.GenerateJWT(user.ID, user.Role)

	cache.Rdb.Set(cache.Ctx, fmt.Sprintf("user:%d", user.ID), token, 24*time.Hour)

	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"error": "failed to generate token",
		})
		return nil
	}
	PublishNotification(user.ID, "candidate", "🎉 You have successfully logged in!")
	c.JSON(&fiber.Map{
		"message": "login successfull",
		"token":   token,
	})

	return nil
}

func GetCandidateProfile(c *fiber.Ctx) error {
	id := c.Locals("user_id")
	candidateModel := &models.Candidate{}
	if id == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := db.DB.Where("id = ?", id).First(candidateModel).Error
	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "could not find the candidate",
		})
		return err
	}

	c.Status(201).JSON(&fiber.Map{
		"message": "candidate found successfull",
		"data":    candidateModel,
	})
	return nil
}

func UpdateCandidateProfile(c *fiber.Ctx) error {
	id := c.Locals("user_id")

	updateUser := &models.UpdateCandidateUser{}
	if err := c.BodyParser(updateUser); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	result := db.DB.Model(&models.Candidate{}).Where("id = ?", id).Updates(updateUser)

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

func ListCandidateApplications(c *fiber.Ctx) error {
	candidateID := c.Locals("user_id") // Set by your middleware
	if candidateID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var applications []models.Application
	if err := db.DB.Preload("Candidate").Where("candidate_id = ?", candidateID).Find(&applications).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Applications not found"})
	}

	if len(applications) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No applications found"})
	}

	var respList []ApplicationResponse

	for _, app := range applications {
		candidate := CandidateSafe{
			Name:      app.Candidate.Name,
			Email:     app.Candidate.Email,
			ResumeUrl: app.Candidate.ResumeUrl,
			Skills:    app.Candidate.Skills,
			Role:      app.Candidate.Role,
		}

		resp := ApplicationResponse{
			Status:    app.Status,
			Candidate: candidate,
			JobID:     app.JobID,
			CreatedAt: app.CreatedAt,
		}

		respList = append(respList, resp)
	}

	return c.JSON(respList)
}

func ApplyForJob(c *fiber.Ctx) error {
	candidateID := c.Locals("user_id")
	jobID := c.Params("id")
	fmt.Println("jobID param:", jobID) // Debug print

	jobIDUint, err := strconv.ParseUint(jobID, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid job ID format",
		})
	}

	if candidateID == nil || jobID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing candidate ID or job ID",
		})
	}

	candidateIDFloat, ok := candidateID.(float64)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid candidate ID type"})
	}
	candidateIDUint := uint(candidateIDFloat)

	// Check if already applied
	var existing models.Application
	if err := db.DB.Where("candidate_id = ? AND job_id = ?", candidateID, jobIDUint).First(&existing).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "Already applied for this job",
		})
	}

	application := models.Application{
		CandidateID: candidateIDUint,
		JobID:       uint(jobIDUint),
		Status:      "applied",
	}

	if err := db.DB.Create(&application).Error; err != nil {
		fmt.Println("Create application error:", err) // Add this line
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create application",
			"details": err.Error(), // Add this for more info
		})
	}
	PublishNotification(candidateIDUint, "candidate", "🎉 You have successfully applied for a job!")
	return c.Status(201).JSON(fiber.Map{
		"message": "Application created successfully",
		"data":    application,
	})
}

func DeleteCandidate(c *fiber.Ctx) error {
	candidateModel := models.Candidate{}
	id := c.Params("id")
	if id == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := db.DB.Delete(candidateModel, id)

	if err.Error != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "could not delete the candidate",
		})
		return err.Error
	}
	PublishNotification(candidateModel.ID, "candidate", "🎉 You have successfully deleted your account!")
	c.Status(200).JSON(&fiber.Map{
		"message": "candidate deleted successull",
	})

	return nil
}
