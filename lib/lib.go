package lib

import (
	"fmt"
	"time"
)

var defaultNotesDir string = "./notes"

func runDailyMigration(notesDir string, current time.Time) error {
	fmt.Println("Daily migration")

	return nil
}

func runMonthlyMigration(notesDir string, currentTime time.Time) error {
	fmt.Println("Monthly migration")

	return nil
}

func RunDailyMigration() error {
	currentTime := time.Now()
	if err := runDailyMigration(defaultNotesDir, currentTime); err != nil {
		return fmt.Errorf("Error running daily migration: %w", err)
	}

	return nil
}

func RunMonthlyMigration() error {
	currentTime := time.Now()
	if err := runMonthlyMigration(defaultNotesDir, currentTime); err != nil {
		return fmt.Errorf("Error running monthly migration: %w", err)
	}

	return nil
}
