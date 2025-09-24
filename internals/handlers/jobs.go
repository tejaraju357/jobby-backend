package handlers

import (
	"jobby/internals/db"
	"jobby/internals/models"

	"github.com/gofiber/fiber/v2"
)

func ListJobs(c *fiber.Ctx) error {
	jobModels := &[]models.Job{}

	err := db.DB.Find(jobModels).Error
	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"error": "could not get the gobs",
		})
		return nil
	}

	c.Status(201).JSON(&fiber.Map{
		"message": "jobs fetched succesfull",
	})
	return nil
}

func GetJobById(c *fiber.Ctx) error {
	id := c.Params("id")
	jobModel := &models.Job{}

	if id == "" {
		c.Status(404).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := db.DB.Where("id = ?", id).First(jobModel).Error
	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "could not get the job with the given id",
		})
		return nil
	}

	c.Status(201).JSON(&fiber.Map{
		"message": "job fetched with id successfully",
	})

	return nil
}

func CreateJob(c *fiber.Ctx) error {
	var job models.Job
	err := c.BodyParser(&job)
	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "invalid request body",
			"details": err.Error(),
		})
		return nil
	}

	if job.Title == "" || job.Description == "" || job.Location == "" || job.Salary == "" || job.CompanyID == 0 {
		c.Status(404).JSON(&fiber.Map{
			"message": "please fill above details",
		})
		return nil
	}

	err = db.DB.Create(&job).Error
	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "failes to register user",
			"details": err.Error(),
		})
		return nil
	}
	PublishNotification(job.CompanyID, "company", "🎉 You have successfully created a job!")
	c.Status(201).JSON(&fiber.Map{
		"message": "job created successfull",
		"job":     job,
	})
	return nil
}

func SearchJobs(c *fiber.Ctx) error {
	title := c.Query("title")
	location := c.Query("location")
	salary := c.Query("salary")

	var jobs []models.Job
	query := db.DB

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}
	if salary != "" {
		query = query.Where("salary = ?", salary)
	}

	err := query.Find(&jobs).Error
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"error": "could not find jobs",
		})
	}

	return c.JSON(&fiber.Map{
		"message": "jobs fetched successfully",
		"jobs":    jobs,
	})
}

func DeleteJob(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		c.Status(405).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	jobModel := &models.Job{}
	errGet := db.DB.Where("id = ?", id).First(jobModel).Error
	if errGet != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "could not find job to delete",
			"details": errGet.Error(),
		})
		return nil
	}

	err := db.DB.Delete(jobModel, id)

	if err != nil {
		c.Status(404).JSON(&fiber.Map{
			"message": "could not deleted job",
		})
		return nil
	}
	PublishNotification(jobModel.CompanyID, "company", "🎉 You have successfully deleted a job!")
	c.Status(201).JSON(&fiber.Map{
		"message": "job deleted successfull",
	})

	return nil
}
