package db

import (
	"jobby/internals/models"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


    var DB *gorm.DB
func DBconnect() {

    dsn := "postgresql://postgres.xnztviszfqpsfumcvdzs:123456789@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres"
    var err error
    DB, err = gorm.Open(postgres.New(postgres.Config{
        DSN:                  dsn,
        PreferSimpleProtocol: true, // 🚀 disables prepared statement caching
    }), &gorm.Config{})
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    // Auto migrate the schema
    DB.AutoMigrate(&models.Application{},&models.Candidate{},&models.Company{},&models.Job{})
}