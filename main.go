package main

import (
	"jobby/internals/db"
	"jobby/internals/handlers"
	"jobby/internals/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db.DBconnect()
	app := fiber.New()

	auth := app.Group("api/")
	auth.Post("/candidate/register", handlers.RegisterCandidate) // done
	auth.Post("/candidate/login", handlers.LoginCandidate)       //done
	auth.Post("/candidate/verify", handlers.VerifyOTPCandidate)  //dome

	auth.Post("/company/register", handlers.RegisterCompany) //dome
	auth.Post("/company/login", handlers.LoginCompany)       //done
	auth.Post("/company/verify", handlers.VerifyOTPCompany)  //done

	routes := app.Group("/api", middleware.Verification)//done
	routes.Get("/candidate/profile", handlers.GetCandidateProfile) //done
	routes.Put("/candidate/profile", handlers.UpdateCandidateProfile)//done
	routes.Get("/candidate/applications", handlers.ListCandidateApplications)//done
	routes.Post("/candidate/apply/:id", handlers.ApplyForJob)//done
	routes.Delete("/candidate/:id", handlers.DeleteCandidate) //done

	routes.Get("/company/profile/:id", handlers.GetCompanyProfile) //done
	routes.Put("/company/profile/:id", handlers.UpdateCompanyProfile) //done
	routes.Get("/company/jobs", handlers.ListCompanyJobs) //done
	routes.Get("/company/jobs/:id/application", handlers.GetApplicationForJob) //done
	routes.Delete("/company/:id", handlers.DeleteCompany) //done

	routes.Get("/applications/:id", handlers.GetApplicationByID)//done
	routes.Get("/applications/candidate/:candidateId", handlers.ListApplicationsByCandidate)//done
	routes.Get("/applications/job/:jobId", handlers.ListApplicationsByJob)//done
	routes.Put("/applications/:id/status", handlers.UpdateApplicationStatus)//done
	routes.Delete("/applications/:id", handlers.DeleteApplication)//done

	routes.Get("/jobs/", handlers.ListJobs)//done
	routes.Get("/jobs/:id", handlers.GetJobById)//done
	routes.Post("/jobs/jobs", handlers.CreateJob)//done
	routes.Get("/jobs/search?title=dev&location=hyd", handlers.SearchJobs)//done
	routes.Delete("/jobs/:id", handlers.DeleteJob)//done


	app.Listen(":8000")

}


