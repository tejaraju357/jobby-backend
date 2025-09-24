package handlers

import (
	"fmt"
	"jobby/internals/cache"
	"jobby/internals/db"
	"jobby/internals/models"
	"jobby/internals/utils"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterCompany(c *fiber.Ctx) error {
	var user models.Company
	err := c.BodyParser(&user)
	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "invalid request body",
			"details": err.Error(),
		})
		return nil
	}

	if user.Name == "" || user.Password == "" || user.Email == "" || user.Description == "" || user.Website == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "please fill above details",
		})
		return nil
	}

	var existing models.Company

	err = db.DB.Where("email = ?", user.Email).First(&existing).Error
	if err == nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "eamil already exixts",
		})
		return err
	}

	err = db.DB.Create(&user).Error
	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "failes to register user",
			"details": err.Error(),
		})
		return nil
	}
	PublishNotification(user.ID, "company", "🎉 You have successfully registered!")
	c.Status(201).JSON(&fiber.Map{
		"message":   "candidate created successfull",
		"candidate": user,
	})
	return nil
}

func LoginCompany(c *fiber.Ctx) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.BodyParser(&credentials)
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

	var user models.Company

	credentials.Email = strings.TrimSpace(strings.ToLower(credentials.Email))
	credentials.Password = strings.TrimSpace(credentials.Password)

	if err := db.DB.Where("email = ? AND password=?", credentials.Email, credentials.Password).First(&user).Error; err != nil {
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

func VerifyOTPCompany(c *fiber.Ctx) error {

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

	var user models.Company
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
	PublishNotification(user.ID, "company", "🎉 You have successfully logged in!")
	c.JSON(&fiber.Map{
		"message": "login successfull",
		"token":   token,
	})

	return nil
}

func GetCompanyProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	companyModel := &models.Company{}
	if id == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := db.DB.Where("id = ?", id).First(companyModel).Error
	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "could not find the company",
		})
		return err
	}

	c.Status(201).JSON(&fiber.Map{
		"message": "compnay found",
		"data":    companyModel,
	})
	return nil
}

func UpdateCompanyProfile(c *fiber.Ctx) error {
	id := c.Params("id")

	updateUser := &models.UpdateCompanyDetails{}
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

func ListCompanyJobs(c *fiber.Ctx) error {
	companyID := c.Locals("user_id")
	if companyID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	var jobs []models.Job
	if err := db.DB.Where("company_id = ?", companyID).Find(&jobs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch jobs"})
	}
	return c.JSON(fiber.Map{"Jobs": jobs})
}

func GetApplicationForJob(c *fiber.Ctx) error {
	jobID := c.Params("job_id")
	if jobID == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	var applications []models.Application
	if err := db.DB.Where("job_id = ?", jobID).Find(&applications).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Applications not found"})
	}
	return c.JSON(fiber.Map{"Applications": applications})
}

func DeleteCompany(c *fiber.Ctx) error {
	companyModel := models.Company{}
	id := c.Params("id")

	if id == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := db.DB.Delete(companyModel, id)

	if err.Error != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "could not delete company",
		})
		return err.Error
	}
	PublishNotification(companyModel.ID, "company", "🎉 You have successfully deleted your account!")
	c.Status(201).JSON(&fiber.Map{
		"message": "Company deleted successfull",
	})

	return nil
}
