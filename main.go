package main

import (
	"log"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Record struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FirstName   string    `json:"first_name"`
	Age         int       `json:"age"`
	DateCreated time.Time `json:"date_created"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("records.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Auto-migrate table
	if err := db.AutoMigrate(&Record{}); err != nil {
		log.Fatal("failed to migrate schema:", err)
	}

	e := echo.New()
	e.POST("/records", createRecord)
	e.GET("/records", getRecords)

	log.Println("Server started at :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func createRecord(c echo.Context) error {
	type Request struct {
		FirstName string `json:"first_name"`
		Age       int    `json:"age"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	record := Record{
		FirstName:   req.FirstName,
		Age:         req.Age,
		DateCreated: time.Now(),
	}

	if err := db.Create(&record).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to insert record"})
	}

	return c.JSON(http.StatusCreated, record)
}

func getRecords(c echo.Context) error {
	var records []Record
	if err := db.Order("id").Find(&records).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch records"})
	}
	return c.JSON(http.StatusOK, records)
}
