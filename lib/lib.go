package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var defaultNotesRootDir string = "./notes"
var defaultTasksFile string = "tasks.note"
var tmpNoteFile string = ".tmp.note"

var errNotesDirDoesNotExist = errors.New("Notes directory does not exist")
var errNextNoteFileExists = errors.New("Next note file already exists")

func monthDir(month time.Month) string {
	return monthPrefix(month)
}

func monthPrefix(month time.Month) string {
	currentMonth := strings.ToLower(month.String())
	return currentMonth[0:3] // Use the first 3 characters, ex. "jul", "aug", "sep", etc.
}

func currentYearDir(currentTime time.Time) string {
	return fmt.Sprintf("%d", currentTime.Year())
}

func previousYearDir(currentTime time.Time) string {
	return fmt.Sprintf("%d", currentTime.Year() - 1)
}

func currentMonthDir(currentTime time.Time) string {
	return monthDir(currentTime.Month())
}

func previousMonthDir(currentTime time.Time) string {
	return monthDir(currentTime.Month() - 1)
}

func nextNoteFile(currentTime time.Time) string {
	return fmt.Sprintf("%s%d.note", monthPrefix(currentTime.Month()), currentTime.Day())
}

func runMigration(notesDir, newFilePath string, ignoredFilePaths []string) error {
	if _, err := os.Stat(newFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to stat notes directory: %w", err)
	} else if os.IsNotExist(err) {
		goto createTargetDirectory
	}

	return errNextNoteFileExists

createTargetDirectory:

	targetMonthDir := filepath.Dir(newFilePath)
	if _, err := os.Stat(targetMonthDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to stat notes directory: %w", err)
	} else if os.IsNotExist(err) {
		// Continue to create target directory
	} else {
		goto migration
	}

	if err := os.MkdirAll(targetMonthDir, 0755); err != nil {
		return fmt.Errorf("Failed to create directory for new note file: %w", err)
	}

migration:

	noteFilePattern := filepath.Join(notesDir, "*.note")
	allNoteFilePaths, err := filepath.Glob(noteFilePattern)
	if err != nil {
		return fmt.Errorf("Failed to find note file paths: %w", err)
	}

	var noteFilePaths []string
	for _, noteFilePath := range(allNoteFilePaths) {
		if slices.ContainsFunc(ignoredFilePaths, func(ignoredFilePath string) bool { return filepath.Base(noteFilePath) == ignoredFilePath }) {
			continue
		}

		noteFilePaths = append(noteFilePaths, noteFilePath)
	}

	noteFileNoteTrees := make(map[string]NoteTree)
	for _, noteFilePath := range(noteFilePaths) {
		noteFileText, err := readFileText(noteFilePath)
		if err != nil {
			return fmt.Errorf("Failed to read note file: %w", err)
		}

		noteTree, err := ParseNoteTree(noteFileText)
		if err != nil {
			return fmt.Errorf("Failed to parse note tree: %w", err)
		}

		noteFileNoteTrees[noteFilePath] = noteTree
	}

	var newNoteTree NoteTree
	for _, noteFilePath := range(noteFilePaths) {
		noteTree := noteFileNoteTrees[noteFilePath]
		noteTreeCopy := noteTree.Copy()
		if err := noteTreeCopy.FilterIncompleteTasks(); err != nil {
			return fmt.Errorf("Failed to filter incomplete tasks: %w", err)
		}

		newNoteTree.Merge(noteTreeCopy)
	}

	if err = os.WriteFile(newFilePath, []byte(newNoteTree.String() + "\n"), 0644); err != nil {
		return fmt.Errorf("Failed to write new note file: %w", err)
	}

	// Use temporary file + rename to ensure that note files are replaced atomically
	tmpNoteFilePath := filepath.Join(notesDir, tmpNoteFile)
	if err = os.RemoveAll(tmpNoteFilePath); err != nil { // Ensure that tmp notes file is removed if it exists
		return fmt.Errorf("Failed to remove temporary notes file: %w", err)
	}

	for _, noteFilePath := range(noteFilePaths) {
		noteTree := noteFileNoteTrees[noteFilePath]
		if err := noteTree.MigrateAll(); err != nil {
			return fmt.Errorf("Failed to migrate notes: %w", err)
		}

		if err = os.WriteFile(tmpNoteFilePath, []byte(noteTree.String() + "\n"), 0644); err != nil {
			return fmt.Errorf("Failed to write new note file: %w", err)
		}

		if err = os.Rename(tmpNoteFilePath, noteFilePath); err != nil {
			return fmt.Errorf("Failed to rename temporary note file: %w", err)
		}
	}

	return nil
}

func runDailyMigration(notesRootDir string, currentTime time.Time) error {
	if _, err := os.Stat(notesRootDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to stat notes directory: %w", err)
	} else if os.IsNotExist(err) {
		return errNotesDirDoesNotExist
	}

	targetMonthDir := filepath.Join(notesRootDir, currentYearDir(currentTime), currentMonthDir(currentTime))
	targetNoteFile := filepath.Join(targetMonthDir, nextNoteFile(currentTime))

	ignoredFilePaths := []string{defaultTasksFile}

	return runMigration(targetMonthDir, targetNoteFile, ignoredFilePaths)
}

func runMonthlyMigration(notesRootDir string, currentTime time.Time) error {
	if _, err := os.Stat(notesRootDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to stat notes directory: %w", err)
	} else if os.IsNotExist(err) {
		return errNotesDirDoesNotExist
	}

	targetMonthDir := filepath.Join(notesRootDir, currentYearDir(currentTime), currentMonthDir(currentTime))
	targetNoteFile := filepath.Join(targetMonthDir, defaultTasksFile)

	var prevNotesDir string
	if currentTime.Month() == time.January {
		prevNotesDir = filepath.Join(notesRootDir, previousYearDir(currentTime), monthDir(time.December))
	} else {
		prevNotesDir = filepath.Join(notesRootDir, currentYearDir(currentTime), previousMonthDir(currentTime))
	}

	ignoredFilePaths := []string{}

	return runMigration(prevNotesDir, targetNoteFile, ignoredFilePaths)
}

func RunDailyMigration() error {
	currentTime := time.Now()
	if err := runDailyMigration(defaultNotesRootDir, currentTime); err != nil {
		return fmt.Errorf("Error running daily migration: %w", err)
	}

	return nil
}

func RunMonthlyMigration() error {
	currentTime := time.Now()
	if err := runMonthlyMigration(defaultNotesRootDir, currentTime); err != nil {
		return fmt.Errorf("Error running monthly migration: %w", err)
	}

	return nil
}
