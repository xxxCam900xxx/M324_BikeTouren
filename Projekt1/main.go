package main

import (
    "database/sql"
    "fmt"
    "log"
	"time"
    "strings"
    "errors"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "myuser"
    password = "mypassword"
    dbname   = "mydb"
)

var db *sql.DB


type Bike struct {
    ID          int       `json:"id"`
    Type        string    `json:"type" binding:"required"`
    FrameNumber string    `json:"frame_number" binding:"required"`
    WheelSize   int       `json:"wheel_size" binding:"required"`
    CreatedAt   time.Time `json:"created_at"`
}

func main() {
	initDB()
}

func initDB() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    var err error
    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal("Could not connect:", err)
    }
    fmt.Println("✅ Connected to Postgres!")

    createTable := `
    CREATE TABLE IF NOT EXISTS bikes (
    id SERIAL PRIMARY KEY,                        
    type VARCHAR(100) NOT NULL,                   
    frame_number VARCHAR(50) NOT NULL UNIQUE,     
    wheel_size INT NOT NULL,                      
    created_at TIMESTAMP NOT NULL DEFAULT NOW()   
	);
	`
    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }
}

func (b *Bike) Validate() error {
    if strings.TrimSpace(b.Type) == "" {
        return errors.New("Bike-Typ darf nicht leer sein")
    }
    if strings.TrimSpace(b.FrameNumber) == "" {
        return errors.New("Rahmennummer darf nicht leer sein")
    }
    if b.WheelSize <= 0 {
        return errors.New("Radhöhe muss größer als 0 sein")
    }
    return nil
}

func InsertBike(db *sql.DB, bike *Bike) error {
    if err := bike.Validate(); err != nil {
        return err
    }

    query := `INSERT INTO bikes (type, frame_number, wheel_size) 
              VALUES ($1, $2, $3) RETURNING id, created_at`

    err := db.QueryRow(query, bike.Type, bike.FrameNumber, bike.WheelSize).
        Scan(&bike.ID, &bike.CreatedAt)

    if err != nil {
        if strings.Contains(err.Error(), "duplicate key value") {
            return errors.New("Rahmennummer existiert bereits")
        }
        return err
    }

    return nil
}

func GetAllBikes(db *sql.DB) ([]Bike, error) {
    query := `
        SELECT id, type, frame_number, wheel_size, created_at
        FROM bikes
        ORDER BY type ASC, wheel_size ASC, created_at ASC;
    `

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var bikes []Bike
    for rows.Next() {
        var b Bike
        if err := rows.Scan(&b.ID, &b.Type, &b.FrameNumber, &b.WheelSize, &b.CreatedAt); err != nil {
            return nil, err
        }
        bikes = append(bikes, b)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return bikes, nil
}