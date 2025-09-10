
# Jobby - Job Application Platform

Jobby is a web application built with **Go (Golang)**, **Fiber framework**, and **GORM ORM** for managing job listings, companies, candidates, and applications.  
It allows companies to post jobs and candidates to apply for them. The platform handles relationships between Jobs, Companies, Candidates, and Applications.

---

## 🚀 Features

- ✅ CRUD operations for Companies, Jobs, Candidates, and Applications  
- ✅ Candidate authentication and authorization  
- ✅ Relationship management between models:
  - A Job belongs to one Company  
  - A Job has many Applications  
  - A Candidate can apply to multiple Jobs  
- ✅ Preloading relations (with GORM) for efficient data retrieval  
- ✅ RESTful API endpoints  

---

## ⚙️ Tech Stack

- Language: Go (Golang)  
- Web Framework: Fiber (https://gofiber.io/)  
- ORM: GORM (https://gorm.io/)  
- Database: PostgreSQL  
- Authentication: Middleware using JWT or session-based (customizable)  





