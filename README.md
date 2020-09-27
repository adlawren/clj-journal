# clj-journal

A digitized bullet journal

Expected directory structure (for example): `notes/2019/dec/*.note`

Unfinished tasks in `.note` files can be automatically migrated. Use `-m` to migrate unfinished tasks from the files in the current month to a new file for the current day. Use `-M` to migrate unfinished tasks from the files in the previous month to a new `tasks.note` file for the current month

See the test data in `resources/test` for concrete examples of notes and the expected directory structure

## Notes

Each note begins with a single character, followed by whitespace, then the note itself. Notes may be indented using whitespace as needed. The leading character may be any of the following:
* '-': A standard note
* '?': A question
* '*': A task
* 'x': A completed task
* '~': A task that no longer needs to be completed, or is invalid
* '>': A task that has been postponed, and moved to a later note
* '<': A task that has been moved to a global task list (ex. monthly tasks, or a more general list of long-term goals)

## Build

```
bash build.sh
```

## Run

```
./clj-journal
```
