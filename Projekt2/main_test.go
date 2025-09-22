package main

import (
    "testing"
    "time"
	"regexp"
    "github.com/DATA-DOG/go-sqlmock"
)

func TestTourValidate(t *testing.T) {
    now := time.Now()

    tests := []struct {
        name    string
        tour    Tour
        wantErr bool
    }{
        {
            name: "valid tour",
            tour: Tour{
                StartLocation: "Bahnhofplatz",
                EndLocation:   "Dorfstrasse",
                StartTime:     now,
                EndTime:       now.Add(2 * time.Hour),
                Companion:     "Sam Meyer",
                Bike:          "Ghost XY1",
            },
            wantErr: false,
        },
        {
            name: "same start and end location",
            tour: Tour{
                StartLocation: "Bahnhofplatz",
                EndLocation:   "Bahnhofplatz",
                StartTime:     now,
                EndTime:       now.Add(2 * time.Hour),
                Companion:     "Sam Meyer",
                Bike:          "Ghost XY1",
            },
            wantErr: true,
        },
        {
            name: "end before start",
            tour: Tour{
                StartLocation: "Bahnhofplatz",
                EndLocation:   "Dorfstrasse",
                StartTime:     now,
                EndTime:       now.Add(-1 * time.Hour),
                Companion:     "Sam Meyer",
                Bike:          "Ghost XY1",
            },
            wantErr: true,
        },
        {
            name: "missing bike",
            tour: Tour{
                StartLocation: "Bahnhofplatz",
                EndLocation:   "Dorfstrasse",
                StartTime:     now,
                EndTime:       now.Add(1 * time.Hour),
                Companion:     "Sam Meyer",
                Bike:          "",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        err := tt.tour.Validate()
        if (err != nil) != tt.wantErr {
            t.Errorf("%s: expected error=%v, got %v", tt.name, tt.wantErr, err)
        }
    }
}

func TestInsertTour(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error initializing sqlmock: %v", err)
    }
    defer db.Close()

    tour := &Tour{
        StartLocation: "Bahnhofplatz",
        EndLocation:   "Dorfstrasse",
        StartTime:     time.Now(),
        EndTime:       time.Now().Add(2 * time.Hour),
        Companion:     "Sam Meyer",
        Bike:          "Ghost XY1",
    }

    mock.ExpectQuery(regexp.QuoteMeta(
        `INSERT INTO tours (start_location, end_location, start_time, end_time, companion, bike)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id, created_at`)).
        WithArgs(tour.StartLocation, tour.EndLocation, tour.StartTime, tour.EndTime, tour.Companion, tour.Bike).
        WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
            AddRow(1, time.Now()))

    err = InsertTour(db, tour)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    if tour.ID != 1 {
        t.Errorf("expected ID=1, got %d", tour.ID)
    }
}

func TestGetAllTours(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error initializing sqlmock: %v", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{
        "id", "start_location", "end_location", "start_time", "end_time", "companion", "bike", "created_at",
    }).AddRow(1, "Bahnhofplatz", "Dorfstrasse", time.Now(), time.Now().Add(2*time.Hour), "Sam Meyer", "Ghost XY1", time.Now())

    mock.ExpectQuery("SELECT id, start_location").
        WillReturnRows(rows)

    tours, err := GetAllTours(db)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    if len(tours) != 1 {
        t.Errorf("expected 1 tour, got %d", len(tours))
    }
}