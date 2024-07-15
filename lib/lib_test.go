package lib

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func copyDir(t *testing.T, dir, dest string) error {
	if err := os.MkdirAll(dest, 0700); err != nil {
		t.Fatalf("Failed copy %s directory to %s: %v", dir, dest, err)
	}

	dirGlobPath := filepath.Join(dir, "*")
	dirPaths, err := filepath.Glob(dirGlobPath)
	if err != nil {
		t.Fatalf("Failed to create directory glob: %v", err)
	}

	for _, filePath := range dirPaths {
		fileContents, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if err = os.WriteFile(filepath.Join(dest, filepath.Base(filePath)), fileContents, 0700); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	return nil
}

func dailyMigrationTime(t *testing.T) time.Time {
	testTime, err := time.Parse(time.DateOnly, "2019-12-25")
	if err != nil {
		t.Fatalf("Failed to parse test time: %v", err)
	}

	return testTime
}

func monthlyMigrationTime(t *testing.T) time.Time {
	testTime, err := time.Parse(time.DateOnly, "2020-01-01")
	if err != nil {
		t.Fatalf("Failed to parse test time: %v", err)
	}

	return testTime
}

func tempNotesDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "bujo_*")
	if err != nil {
		t.Fatalf("Failed to create temporary test directory: %v", err)
	}

	notesDir := filepath.Join(tempDir, "notes")
	if err := os.Mkdir(notesDir, 0700); err != nil {
		t.Fatalf("Failed to create temporary notes directory: %v", err)
	}

	return notesDir
}

func testFilesEqual(t *testing.T, dir1, dir2 string) bool {
	dir1GlobPath := filepath.Join(dir1, "*")
	dir2GlobPath := filepath.Join(dir2, "*")

	dir1Paths, err := filepath.Glob(dir1GlobPath)
	if err != nil {
		t.Fatalf("Failed to create directory glob: %v", err)
	}
	dir2Paths, err := filepath.Glob(dir2GlobPath)
	if err != nil {
		t.Fatalf("Failed to create directory glob: %v", err)
	}

	if len(dir1Paths) != len(dir2Paths) {
		return false
	}

	// Compare files
	for i := range dir1Paths {
		dir1Path := dir1Paths[i]
		dir2Path := dir2Paths[i]

		dir1PathBasename := filepath.Base(dir1Path)
		dir2PathBasename := filepath.Base(dir2Path)
		if dir1PathBasename != dir2PathBasename {
			return false
		}

		dir1FileContents, err := os.ReadFile(dir1Path)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		dir2FileContents, err := os.ReadFile(dir2Path)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if !bytes.Equal(dir1FileContents, dir2FileContents) {
			return false
		}
	}

	return true
}

func TestRunDailyMigration(t *testing.T) {
	notesDir := tempNotesDir(t)
	t.Logf("Using temporary notes directory: %s", notesDir)

	copyDir(t, "./test/dec", filepath.Join(notesDir, "2019", "dec"))

	if err := runDailyMigration(notesDir, dailyMigrationTime(t)); err != nil {
		t.Fatalf("Failed to run daily migration: %v", err)
	}

	if !testFilesEqual(t, "./test/expected-dec", filepath.Join(notesDir, "2019", "dec")) {
		t.Fatal("Migrated files do not match expected files")
	}
}

func TestRunDailyMigrationReturnsErrorIfNotesDirectoryDoesNotExist(t *testing.T) {
	if err := runDailyMigration("non-existent-dir", dailyMigrationTime(t)); !errors.Is(err, errNotesDirDoesNotExist) {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestRunDailyMigrationReturnsErrorIfTargetFileAlreadyExists(t *testing.T) {
	notesDir := tempNotesDir(t)
	t.Logf("Using temporary notes directory: %s", notesDir)

	copyDir(t, "./test/dec", filepath.Join(notesDir, "2019", "dec"))

	if err := os.WriteFile(filepath.Join(notesDir, "2019", "dec", "dec25.note"), []byte("- This is a test"), 0700); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	if err := runDailyMigration(notesDir, dailyMigrationTime(t)); !errors.Is(err, errNextNoteFileExists) {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestRunMonthlyMigration(t *testing.T) {
	notesDir := tempNotesDir(t)
	t.Logf("Using temporary notes directory: %s", notesDir)

	copyDir(t, "./test/dec", filepath.Join(notesDir, "2019", "dec"))

	if err := runMonthlyMigration(notesDir, monthlyMigrationTime(t)); err != nil {
		t.Fatalf("Failed to run monthly migration: %v", err)
	}

	if !testFilesEqual(t, "./test/expected-jan", filepath.Join(notesDir, "2020", "jan")) {
		t.Fatal("Migrated files do not match expected files")
	}
}

func TestRunMonthlyMigrationReturnsErrorIfNotesDirectoryDoesNotExist(t *testing.T) {
	if err := runMonthlyMigration("non-existent-dir", monthlyMigrationTime(t)); !errors.Is(err, errNotesDirDoesNotExist) {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestRunMonthlyMigrationReturnsErrorIfTargetDirectoryAlreadyExists(t *testing.T) {
	notesDir := tempNotesDir(t)
	t.Logf("Using temporary notes directory: %s", notesDir)

	copyDir(t, "./test/dec", filepath.Join(notesDir, "2019", "dec"))

	if err := os.MkdirAll(filepath.Join(notesDir, "2020", "jan"), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	if err := runMonthlyMigration(notesDir, monthlyMigrationTime(t)); !errors.Is(err, errNextMonthDirExists) {
		t.Fatalf("Unexpected error: %v", err)
	}
}
