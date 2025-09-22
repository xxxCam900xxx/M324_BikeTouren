package main

import (
    "errors"
    "strings"
    "time"
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"

)

func main() {
	InitDB()
}

func InitDB() *sql.DB {
    connStr := "host=localhost port=5433 user=myuser password=mypassword dbname=mydb sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal("Could not connect:", err)
    }

    createTable := `
    CREATE TABLE IF NOT EXISTS tours (
        id SERIAL PRIMARY KEY,
        start_location VARCHAR(255) NOT NULL,
        end_location VARCHAR(255) NOT NULL,
        start_time TIMESTAMP NOT NULL,
        end_time TIMESTAMP NOT NULL,
        companion VARCHAR(100) NOT NULL,
        bike VARCHAR(100) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );`

    if _, err := db.Exec(createTable); err != nil {
        log.Fatal(err)
    }

    fmt.Println("✅ Connected & Table ready")
    return db
}

type Tour struct {
    ID            int       `json:"id"`
    StartLocation string    `json:"start_location" binding:"required"`
    EndLocation   string    `json:"end_location" binding:"required"`
    StartTime     time.Time `json:"start_time" binding:"required"`
    EndTime       time.Time `json:"end_time" binding:"required"`
    Companion     string    `json:"companion" binding:"required"`
    Bike          string    `json:"bike" binding:"required"`
    CreatedAt     time.Time `json:"created_at"`
}

// Validate prüft die Eingaben
func (t *Tour) Validate() error {
    if strings.TrimSpace(t.StartLocation) == "" {
        return errors.New("Startort darf nicht leer sein")
    }
    if strings.TrimSpace(t.EndLocation) == "" {
        return errors.New("Ankunftsort darf nicht leer sein")
    }
    if t.StartLocation == t.EndLocation {
        return errors.New("Start- und Ankunftsort dürfen nicht identisch sein")
    }
    if t.EndTime.Before(t.StartTime) {
        return errors.New("Endzeit muss nach Startzeit liegen")
    }
    if strings.TrimSpace(t.Companion) == "" {
        return errors.New("Begleiter darf nicht leer sein")
    }
    if strings.TrimSpace(t.Bike) == "" {
        return errors.New("Bike darf nicht leer sein")
    }
    return nil
}

func InsertTour(db *sql.DB, tour *Tour) error {
    if err := tour.Validate(); err != nil {
        return err
    }

    query := `INSERT INTO tours (start_location, end_location, start_time, end_time, companion, bike)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id, created_at`

    return db.QueryRow(query,
        tour.StartLocation,
        tour.EndLocation,
        tour.StartTime,
        tour.EndTime,
        tour.Companion,
        tour.Bike,
    ).Scan(&tour.ID, &tour.CreatedAt)
}

func GetAllTours(db *sql.DB) ([]Tour, error) {
    query := `
        SELECT id, start_location, end_location, start_time, end_time, companion, bike, created_at
        FROM tours
        ORDER BY start_time ASC;
    `

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tours []Tour
    for rows.Next() {
        var t Tour
        if err := rows.Scan(&t.ID, &t.StartLocation, &t.EndLocation, &t.StartTime, &t.EndTime, &t.Companion, &t.Bike, &t.CreatedAt); err != nil {
            return nil, err
        }
        tours = append(tours, t)
    }

    return tours, rows.Err()
}