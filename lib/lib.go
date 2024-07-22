package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var defaultNotesDir string = "./notes"

var errNotesDirDoesNotExist = errors.New("Notes directory does not exist")
var errNextMonthDirExists = errors.New("Next month directory already exists")
var errNextNoteFileExists = errors.New("Next note file already exists")

func monthPrefix(month time.Month) string {
	currentMonth := strings.ToLower(month.String())
	return currentMonth[0:3] // Use the first 3 characters, ex. "jul", "aug", "sep", etc.
}

func currentYearDir(currentTime time.Time) string {
	return fmt.Sprintf("%d", currentTime.Year())
}

func currentMonthDir(currentTime time.Time) string {
	return monthPrefix(currentTime.Month())
}

func nextNoteFile(currentTime time.Time) string {
	return fmt.Sprintf("%s%d.note", monthPrefix(currentTime.Month()), currentTime.Day())
}

func runDailyMigration(notesDir string, currentTime time.Time) error {
	// Advance declarations before goto
	var targetNoteFile string

	if _, err := os.Stat(notesDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to stat notes directory: %w", err)
	} else if os.IsNotExist(err) {
		return errNotesDirDoesNotExist
	}

	targetMonthDir := filepath.Join(notesDir, currentYearDir(currentTime), currentMonthDir(currentTime))
	targetNoteFile = filepath.Join(targetMonthDir, nextNoteFile(currentTime))
	if _, err := os.Stat(targetNoteFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to stat notes directory: %w", err)
	} else if os.IsNotExist(err) {
		goto dailyMigration
	}

	return errNextNoteFileExists

dailyMigration:

	return nil
}

func runMonthlyMigration(notesDir string, currentTime time.Time) error {
	// Advance declarations before goto
	var targetMonthDir string

	if _, err := os.Stat(notesDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to stat notes directory: %w", err)
	} else if os.IsNotExist(err) {
		return errNotesDirDoesNotExist
	}

	targetMonthDir = filepath.Join(notesDir, currentYearDir(currentTime), currentMonthDir(currentTime))
	if _, err := os.Stat(targetMonthDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to stat notes directory: %w", err)
	} else if os.IsNotExist(err) {
		goto monthlyMigration
	}

	return errNextMonthDirExists

monthlyMigration:

	// TODO: Fix permissions, this is required for the tests in the tmp directory, but it's not ideal for everyday use
	if err := os.MkdirAll(targetMonthDir, 0700); err != nil {
		return fmt.Errorf("Failed to create next month directory: %w", err)
	}

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
