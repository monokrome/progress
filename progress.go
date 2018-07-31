package progress

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Tag provides information describing a specific tag
type Tag struct {
	Identifier int
	Name       string
	References []interface{}
}

// Project provides information describing a specific project
type Project struct {
	Identifier   int
	Abbreviation string
	Name         string
	Description  string
}

// Task provides information describing a specific task
type Task struct {
	Identifier int
	Project    *Project
	Summary    string
	Details    string
}

// Database represents a single database containing project data
type Database struct {
	path       string
	connection *sql.DB
}

// Tasks gets a list of recent tasks
func (project *Project) Tasks() (tasks []Task) {
	return tasks
}

func (database *Database) autoUpdateField(table string) error {
	statement, err := database.connection.Prepare(
		"CREATE TRIGGER IF NOT EXISTS " + table + "_auto_update_trigger " + " AFTER UPDATE ON " + table + `
		BEGIN
			UPDATE ` + table + ` SET last_updated = datetime('now')
			WHERE id = NEW.id;
		END;
	`)

	if err != nil {
		return err
	}

	_, err = statement.Exec()
	return err
}

func (database *Database) ensureTable(name string, description string) error {
	statement, err := database.connection.Prepare(
		"CREATE TABLE IF NOT EXISTS " + name + " (" + `
			id           INTEGER PRIMARY KEY AUTOINCREMENT, ` +
			description + `
			created      TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		)`,
	)

	if err != nil {
		return err
	}

	_, err = statement.Exec()

	if err != nil {
		return err
	}

	return database.autoUpdateField(name)
}

func (database *Database) ensureProjectsTableExists() error {
	return database.ensureTable("projects", `
		abbreviation TEXT NOT NULL UNIQUE,
		name         TEXT NOT NULL UNIQUE,
		description  TEXT,
	`)
}

func (database *Database) ensureTasksTableExists() error {
	return database.ensureTable("tasks", `
		project_id   TEXT NOT NULL,
		summary      TEXT NOT NULL,
		description  TEXT,
	`)
}

func (database *Database) ensureTagsTableExists() error {
	err := database.ensureTable("tags", `
		name TEXT NOT NULL,
	`)

	if err != nil {
		return err
	}

	// TODO: Associative data
	return database.ensureTable("tasks_tags", `
		tag_id TEXT NOT NULL,
		task_id TEXT NOT NULL,
	`)
}

func (database *Database) prepare() error {
	if err := database.ensureProjectsTableExists(); err != nil {
		return err
	}

	if err := database.ensureTasksTableExists(); err != nil {
		return err
	}

	if err := database.ensureTagsTableExists(); err != nil {
		return err
	}

	return nil
}

// Open connects to a database and returns a reference to it
func Open(driver string, path string) (*Database, error) {
	var database Database

	if _, err := os.Stat(path); os.IsNotExist(err) {
		directory := filepath.Dir(path)

		if _, err = os.Stat(directory); os.IsNotExist(err) {
			err = os.MkdirAll(directory, 0750)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not create data directory at %v: %v\n", directory, err)
				os.Exit(100)
			}
		}

		if _, err = os.Stat(path); os.IsNotExist(err) {
			if _, err = os.Create(path); err != nil {
				fmt.Fprintf(os.Stderr, "Could not create data file at %v: %v\n", path, err)
				os.Exit(101)
			}
		}
	}

	connection, err := sql.Open(driver, path)

	if err != nil {
		return &database, err
	}

	database.path = path
	database.connection = connection

	if err := database.prepare(); err != nil {
		return nil, err
	}

	return &database, nil
}

// Projects returns a slice of Projects in the database
func (database *Database) Projects() ([]Project, error) {
	var projects []Project

	results, err := database.connection.Query(`
		SELECT
			id,
			abbreviation,
			name,
			description
		FROM projects
	`)

	if err != nil {
		return nil, err
	}

	for results.Next() {
		project := Project{}
		results.Scan(&project.Identifier, &project.Abbreviation, &project.Name, &project.Description)
		projects = append(projects, project)
	}

	return projects, nil
}

// AddProject inserts a new project into the database
func (database *Database) AddProject(name string, abbreviation string, description string) error {
	statement, err := database.connection.Prepare(`
		INSERT INTO projects (name, abbreviation, description) VALUES (?, ?, ?)
	`)

	if err != nil {
		return err
	}

	_, err = statement.Exec(name, abbreviation, description)

	if err != nil {
		return err
	}

	return nil
}

// DefaultProject returns the default project when one is not provided
func (database *Database) DefaultProject() (Project, error) {
	var count int

	project := Project{}

	result, err := database.connection.Query(`
		SELECT
			id,
			name,
			abbreviation,
			description
		FROM projects
		ORDER BY last_updated
		DESC LIMIT 1
	`)

	if err != nil {
		return project, err
	}

	for result.Next() {
		result.Scan(&project.Identifier, &project.Name, &project.Abbreviation, &project.Description)
		count++
	}

	if count < 1 {
		return project, errors.New("no projects exist - have you created one?")
	}

	if count > 1 {
		fmt.Fprintf(os.Stderr, "Expected 1 result when querying default project, but got %v.\n", count)
	}

	return project, nil
}

// Project returns a project with the given name or abbreviation
func (database *Database) Project(reference string) (Project, error) {
	var project Project

	result, err := database.connection.Query(`
		SELECT
			id,
			name,
			abbreviation,
			description
		FROM projects
		WHERE name == ? OR abbreviation == ?
		LIMIT 1
	`, reference, reference)

	if err != nil {
		return project, err
	}

	for result.Next() {
		result.Scan(&project.Identifier, &project.Name, &project.Abbreviation, &project.Description)
	}

	return project, nil
}

// AddTask adds a task to the provided project
func (database *Database) AddTask(project Project, summary string) (int64, error) {
	statement, err := database.connection.Prepare(`
		INSERT INTO tasks (project_id, summary) VALUES (?, ?)
	`)

	if err != nil {
		return -1, err
	}

	result, err := statement.Exec(project.Identifier, summary)

	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
