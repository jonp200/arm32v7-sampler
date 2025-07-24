package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type Record struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	Age         int       `json:"age"`
	DateCreated time.Time `json:"date_created"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./records.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	_, err = db.Exec(
		`
		CREATE TABLE IF NOT EXISTS records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT,
			age INTEGER,
			date_created DATETIME
		);
	`,
	)
	if err != nil {
		log.Fatal(err)
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

	now := time.Now()

	result, err := db.Exec(
		`INSERT INTO records (first_name, age, date_created) VALUES (?, ?, ?)`,
		req.FirstName, req.Age, now,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to insert record"})
	}

	id, _ := result.LastInsertId()
	return c.JSON(
		http.StatusCreated, echo.Map{
			"id":           id,
			"first_name":   req.FirstName,
			"age":          req.Age,
			"date_created": now,
		},
	)
}

func getRecords(c echo.Context) error {
	rows, err := db.Query(`SELECT id, first_name, age, date_created FROM records`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch records"})
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.FirstName, &r.Age, &r.DateCreated); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Scan error"})
		}
		records = append(records, r)
	}

	return c.JSON(http.StatusOK, records)
}
