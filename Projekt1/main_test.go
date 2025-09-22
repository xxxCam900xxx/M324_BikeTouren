package main

import (
    "testing"
    "time"
	"regexp"
    "github.com/DATA-DOG/go-sqlmock"
)

func TestBikeValidate(t *testing.T) {
    tests := []struct {
        name    string
        bike    Bike
        wantErr bool
    }{
        {
            name: "valid bike",
            bike: Bike{Type: "Ghost XY1", FrameNumber: "XL123", WheelSize: 15},
            wantErr: false,
        },
        {
            name: "missing type",
            bike: Bike{Type: "", FrameNumber: "XL123", WheelSize: 15},
            wantErr: true,
        },
        {
            name: "missing frame number",
            bike: Bike{Type: "Ghost XY1", FrameNumber: "", WheelSize: 15},
            wantErr: true,
        },
        {
            name: "invalid wheel size",
            bike: Bike{Type: "Ghost XY1", FrameNumber: "XL123", WheelSize: 0},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        err := tt.bike.Validate()
        if (err != nil) != tt.wantErr {
            t.Errorf("%s: expected error=%v, got %v", tt.name, tt.wantErr, err)
        }
    }
}

func TestInsertBike(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error initializing sqlmock: %v", err)
    }
    defer db.Close()

    bike := &Bike{
        Type:        "Ghost XY1",
        FrameNumber: "XL123",
        WheelSize:   15,
    }

    // Erwartetes SQL
    mock.ExpectQuery(regexp.QuoteMeta(
        `INSERT INTO bikes (type, frame_number, wheel_size) 
              VALUES ($1, $2, $3) RETURNING id, created_at`)).
        WithArgs(bike.Type, bike.FrameNumber, bike.WheelSize).
        WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
            AddRow(1, time.Now()))

    err = InsertBike(db, bike)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    if bike.ID != 1 {
        t.Errorf("expected ID=1, got %d", bike.ID)
    }
}

func TestGetAllBikes(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error initializing sqlmock: %v", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{
        "id", "type", "frame_number", "wheel_size", "created_at",
    }).AddRow(1, "Ghost XY1", "XL123", 15, time.Now())

    mock.ExpectQuery("SELECT id, type, frame_number").
        WillReturnRows(rows)

    bikes, err := GetAllBikes(db)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    if len(bikes) != 1 {
        t.Errorf("expected 1 bike, got %d", len(bikes))
    }
    if bikes[0].Type != "Ghost XY1" {
        t.Errorf("expected type Ghost XY1, got %s", bikes[0].Type)
    }
}
