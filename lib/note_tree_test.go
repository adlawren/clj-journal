package lib

import "testing"

func readNoteFile(t *testing.T) string {
	noteFileText, err := readFileText("./test/test.note")
	if err != nil {
		t.Fatalf("Failed to read note file: %v", err)
	}

	return noteFileText
}

func TestFilterIncompleteTasks(t *testing.T) {
	noteFileText := readNoteFile(t)

	noteTree, err := ParseNoteTree(noteFileText)
	if err != nil {
		t.Fatalf("Failed to parse note tree: %v", err)
	}

	if err := noteTree.FilterIncompleteTasks(); err != nil {
		t.Fatalf("Failed to filter incomplete tasks: %v", err)
	}

	if noteCount := len(noteTree.Notes); noteCount != 2 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[0].Text; noteText != "- Test 1\n\ncontent 1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[1].Text; noteText != "- Test 2\n\ncontent 2\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[0].ChildNotes.Notes); noteCount != 2 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[0].ChildNotes.Notes[0].Text; noteText != "  - Test 1.1\n\ncontent 1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[0].ChildNotes.Notes[1].Text; noteText != "  * Test 1.2\n\ncontent 1.2\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes); noteCount != 1 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes[0].Text; noteText != "    * Test 1.1.1\n\ncontent 1.1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[1].ChildNotes.Notes); noteCount != 1 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].Text; noteText != " - Test 2.1\n\ncontent 2.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes); noteCount != 1 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes[0].Text; noteText != "  * Test 2.1.2\n\ncontent 2.1.2\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}
}

func TestMigrateAll(t *testing.T) {
	noteFileText := readNoteFile(t)

	noteTree, err := ParseNoteTree(noteFileText)
	if err != nil {
		t.Fatalf("Failed to parse note tree: %v", err)
	}

	if err := noteTree.MigrateAll(); err != nil {
		t.Fatalf("Failed to migrate notes: %v", err)
	}

	if noteCount := len(noteTree.Notes); noteCount != 3 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[0].Text; noteText != "- Test 1\n\ncontent 1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[1].Text; noteText != "- Test 2\n\ncontent 2\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[2].Text; noteText != "- Test 3\n\ncontent 3\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[0].ChildNotes.Notes); noteCount != 2 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[0].ChildNotes.Notes[0].Text; noteText != "  - Test 1.1\n\ncontent 1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[0].ChildNotes.Notes[1].Text; noteText != "  > Test 1.2\n\ncontent 1.2\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes); noteCount != 1 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes[0].Text; noteText != "    > Test 1.1.1\n\ncontent 1.1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[1].ChildNotes.Notes); noteCount != 1 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].Text; noteText != " - Test 2.1\n\ncontent 2.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes); noteCount != 2 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes[0].Text; noteText != "  - Test 2.1.1\n\ncontent 2.1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes[1].Text; noteText != "  > Test 2.1.2\n\ncontent 2.1.2\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes); noteCount != 1 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes[0].Text; noteText != "    - Test 2.1.1.1\n\ncontent 2.1.1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}
}

func TestParseNoteTree(t *testing.T) {
	noteFileText := readNoteFile(t)

	noteTree, err := ParseNoteTree(noteFileText)
	if err != nil {
		t.Fatalf("Failed to parse note tree: %v", err)
	}

	if noteCount := len(noteTree.Notes); noteCount != 3 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[0].Text; noteText != "- Test 1\n\ncontent 1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[1].Text; noteText != "- Test 2\n\ncontent 2\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[2].Text; noteText != "- Test 3\n\ncontent 3\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[0].ChildNotes.Notes); noteCount != 2 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[0].ChildNotes.Notes[0].Text; noteText != "  - Test 1.1\n\ncontent 1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[0].ChildNotes.Notes[1].Text; noteText != "  * Test 1.2\n\ncontent 1.2\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes); noteCount != 1 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes[0].Text; noteText != "    * Test 1.1.1\n\ncontent 1.1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[1].ChildNotes.Notes); noteCount != 1 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].Text; noteText != " - Test 2.1\n\ncontent 2.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes); noteCount != 2 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes[0].Text; noteText != "  - Test 2.1.1\n\ncontent 2.1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes[1].Text; noteText != "  * Test 2.1.2\n\ncontent 2.1.2\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}

	if noteCount := len(noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes); noteCount != 1 {
		t.Fatalf("Incorrect note count in tree: %d", noteCount)
	}

	if noteText := noteTree.Notes[1].ChildNotes.Notes[0].ChildNotes.Notes[0].ChildNotes.Notes[0].Text; noteText != "    - Test 2.1.1.1\n\ncontent 2.1.1.1\n" {
		t.Fatalf("Unexpected note: %s", noteText)
	}
}
