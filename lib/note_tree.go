package lib

import (
	"fmt"
	"regexp"
	"strings"
)

var leadingWhitespaceRegexString string = "^\\s*"
var bulletRegexString string = "[?\\-*x~><]"

type Note struct {
	Depth int
	Text string
	ChildNotes NoteTree
}

type NoteTree struct {
	Notes []*Note
}

type Stack struct {
	Items []interface{}
}

func (n *Note) AddText(text string) {
	n.Text = n.Text + "\n" + text
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

func (noteTree NoteTree) Length() int {
	return len(noteTree.Notes)
}

func (noteTree NoteTree) String() string {
	var noteStrings []string
	for _, note := range(noteTree.Notes) {
		noteStrings = append(noteStrings, note.String())
	}

	return strings.Join(noteStrings, "\n")
}

func (s Stack) Peek() interface{} {
	if len(s.Items) == 0 {
		return nil
	}

	lastItem := s.Items[len(s.Items)-1]
	return lastItem
}

func (s *Stack) Pop() interface{} {
	if len(s.Items) == 0 {
		return nil
	}

	lastItem := s.Items[len(s.Items)-1]
	s.Items = s.Items[:len(s.Items)-1]
	return lastItem
}

func (s *Stack) Push(item interface{}) {
	s.Items = append(s.Items, item)
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
	for _, line := range(lines) {
		r, err := regexp.Compile(leadingWhitespaceRegexString + bulletRegexString)
		if err != nil {
			return notes, fmt.Errorf("Failed to compile regex: %w", err)
		}

		if !r.MatchString(line) { // If line isn't a bullet note, append it to the previous note
			prevNote := notes[len(notes)-1]
			prevNote.AddText(line)
			continue
		}

		depth, err := lineDepth(line)
		if err != nil {
			return notes, fmt.Errorf("Failed to get line depth: %w", err)
		}

		notes = append(notes, &Note{Text: line, Depth: depth})
	}

	return notes, nil
}


func parseNoteTrees(notes []*Note) (NoteTree, error) {
	var noteStack Stack

	rootNote := Note{Depth: -1}
	noteStack.Push(&rootNote)

	var prevNote *Note
	for _, note := range(notes) {
		var currentDepth int
		if prevNote == nil {
			currentDepth = 0
		} else {
			currentDepth = prevNote.Depth
		}

		parentNote := noteStack.Peek().(*Note)

		if note.Depth < currentDepth {
			// Pop parent notes until we reach the right depth
			for parentNote.Depth >= note.Depth {
				parentNote = noteStack.Pop().(*Note)
			}
		} else if note.Depth > currentDepth {
			parentNote = prevNote
			noteStack.Push(prevNote)
		}

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
