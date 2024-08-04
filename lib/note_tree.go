package lib

import (
	"fmt"
	"regexp"
	"strings"
)

var leadingWhitespaceRegexString string = "^\\s*"
var bulletRegexString string = "[?\\-*x~><]"
var unmigratedBulletRegex string = "\\*"

type Note struct {
	Text string
	Depth int
	ChildNotes NoteTree
}

type NoteTree struct {
	Notes []*Note
}

func (n *Note) AddText(text string) {
	n.Text = n.Text + "\n" + text
}

func (n Note) IsUnmigrated() (bool, error) {
	r, err := regexp.Compile(leadingWhitespaceRegexString + unmigratedBulletRegex)
	if err != nil {
		return false, fmt.Errorf("Failed to compile regex: %w", err)
	}

	return r.MatchString(n.Text), nil
}

func (n *Note) Migrate() error {
	r, err := regexp.Compile(leadingWhitespaceRegexString + unmigratedBulletRegex)
	if err != nil {
		return fmt.Errorf("Failed to compile regex: %w", err)
	}

	leadingBulletText := r.FindString(n.Text)
	newLeadingBulletText := leadingBulletText[:len(leadingBulletText)-1] + ">" // Change the bullet from "*" to ">"
	n.Text = strings.Replace(n.Text, leadingBulletText, newLeadingBulletText, 1)

	return nil
}

func (n Note) String() string {
	noteString := n.Text
	if n.ChildNotes.Length() > 0 {
		return noteString + "\n" + n.ChildNotes.String()
	}

	return noteString
}

func (noteTree *NoteTree) Add(n *Note) {
	noteTree.Notes = append(noteTree.Notes, n)
}

func (noteTree NoteTree) Copy() NoteTree {
	var notes []*Note
	for _, existingNote := range(noteTree.Notes) {
		newChildNotes := existingNote.ChildNotes.Copy()
		notes = append(notes, &Note{Text: existingNote.Text, Depth: existingNote.Depth, ChildNotes: newChildNotes})
	}

	return NoteTree{Notes: notes}
}

func (noteTree *NoteTree) FilterIncompleteTasks() error {
	if noteTree.Length() == 0 {
		return nil
	}

	var newNotes []*Note
	for _, note := range(noteTree.Notes) {
		isUnmigrated, err := note.IsUnmigrated()
		if err != nil {
			return fmt.Errorf("Failed to check if note is unmigrated: %w", err)
		}
		if isUnmigrated {
			newNotes = append(newNotes, note)
			continue
		}

		if err := note.ChildNotes.FilterIncompleteTasks(); err != nil {
			return fmt.Errorf("Failed to filter unmigrated child notes: %w", err)
		}
		if note.ChildNotes.Length() > 0 {
			newNotes = append(newNotes, note)
		}
	}

	noteTree.Notes = newNotes

	return nil
}

func (noteTree NoteTree) Length() int {
	return len(noteTree.Notes)
}

func (noteTree *NoteTree) Merge(otherNoteTree NoteTree) {
	for _, note := range(otherNoteTree.Notes) {
		noteTree.Notes = append(noteTree.Notes, note)
	}
}

func (noteTree NoteTree) MigrateAll() error {
	for _, note := range(noteTree.Notes) {
		isUnmigrated, err := note.IsUnmigrated()
		if err != nil {
			return fmt.Errorf("Failed to check if note is unmigrated: %w", err)
		}
		if !isUnmigrated {
			goto migrateChildNotes
		}
		if err := note.Migrate(); err != nil {
			return fmt.Errorf("Failed to migrate note: %w", err)
		}

	migrateChildNotes:

		if err := note.ChildNotes.MigrateAll(); err != nil {
			return fmt.Errorf("Failed to migrate child notes: %w", err)
		}
	}

	return nil
}

func (noteTree NoteTree) String() string {
	var noteStrings []string
	for _, note := range(noteTree.Notes) {
		noteStrings = append(noteStrings, note.String())
	}

	return strings.Join(noteStrings, "\n")
}

func lineDepth(line string) (int, error) {
	r, err := regexp.Compile(leadingWhitespaceRegexString)
	if err != nil {
		return -1, fmt.Errorf("Failed to compile regex: %w", err)
	}

	matches := r.FindAllString(line, 1)
	if len(matches) == 0 {
		return 0, nil
	}

	return len(matches[0]), nil
}

func parseNotes(text string) ([]*Note, error) {
	lines := strings.Split(text, "\n")

	var notes []*Note
	var prevNote *Note // Advance declaration for goto
	for _, line := range(lines) {
		r, err := regexp.Compile(leadingWhitespaceRegexString + bulletRegexString)
		if err != nil {
			return notes, fmt.Errorf("Failed to compile regex: %w", err)
		}

		if r.MatchString(line) {
			goto parseBulletNote
		}

		// Append non-bullet notes to the previous note
		if len(notes) == 0 {
			continue // Skip non-bullet notes at the top of the file
		}
		prevNote = notes[len(notes)-1]
		prevNote.AddText(line)

		continue

	parseBulletNote:

		depth, err := lineDepth(line)
		if err != nil {
			return notes, fmt.Errorf("Failed to get line depth: %w", err)
		}

		notes = append(notes, &Note{Text: line, Depth: depth})
	}

	return notes, nil
}

// Note: Implementing this recursively results in a less readable implementation
// You need to keep track of the current index in the notes array throughout the recursive calls
// You can pass the index by reference, but it makes the code much more complicated
func parseNoteTrees(notes []*Note) (NoteTree, error) {
	var noteStack Stack

	rootNote := Note{Depth: -1}
	noteStack.Push(&rootNote)

	var currentDepth int
	var prevNote *Note
	for _, note := range(notes) {
		parentNote := noteStack.Peek().(*Note)

		if note.Depth < currentDepth {
			// Pop parent notes until we reach the right depth
			for parentNote.Depth >= note.Depth {
				noteStack.Pop()
				parentNote = noteStack.Peek().(*Note)
			}
		} else if note.Depth > currentDepth {
			parentNote = prevNote
			noteStack.Push(prevNote)
		}

		currentDepth = note.Depth
		prevNote = note
		parentNote.ChildNotes.Add(note)
	}

	return rootNote.ChildNotes, nil
}

func ParseNoteTree(text string) (NoteTree, error) {
	var noteTree NoteTree

	notes, err := parseNotes(text)
	if err != nil {
		return noteTree, fmt.Errorf("Failed to parse notes: %w", err)
	}

	if noteTree, err = parseNoteTrees(notes); err != nil {
		return noteTree, fmt.Errorf("Failed to parse note trees: %w", err)
	} else {
		return noteTree, nil
	}
}
