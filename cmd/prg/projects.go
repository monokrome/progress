package main

import (
	"flag"
	"fmt"
	"github.com/monokrome/progress"
	"os"
	"strings"
)

func getAbbreviation(name string) string {
	if len(name) < 3 {
		return strings.ToUpper(name)
	}

	// This seems to work okay
	return strings.ToUpper(string(name[0]) + name[2:4])
}

func getDatabase() *progress.Database {
	database, err := progress.Open("sqlite3", projectFilePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open data file: %v\n", err)
		os.Exit(2)
	}

	return database
}

func projectListCommand(arguments ...string) {
	if len(arguments) > 0 {
		fmt.Fprintf(os.Stderr, "usage: @project list\n")
		flag.Usage()
		os.Exit(1)
	}

	projects, err := getDatabase().Projects()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project list: %v\n", err)
		os.Exit(3)
	}

	for _, project := range projects {
		fmt.Printf("[%v] %v - %v\n", project.Abbreviation, project.Name, project.Description)
	}
}

func projectAddCommand(arguments ...string) {
	var project string
	var abbreviation string
	var description string

	if 1 > len(arguments) {
		fmt.Fprintf(os.Stderr, "usage: @project add <project name> [abbreviation] [...description]\n")
		os.Exit(1)
	}

	project = arguments[0]

	if len(arguments) > 1 && arguments[1] != "_" {
		abbreviation = arguments[1]
	} else {
		abbreviation = getAbbreviation(project)
	}

	if len(arguments) > 2 {
		description = strings.Join(arguments[2:], " ")
	}

	if len(abbreviation) >= len(project) {
		fmt.Fprintf(os.Stderr, "Error: %v (abbreviation) is not shorter than %v (project name)\n", abbreviation, project)
		os.Exit(5)
	}

	fmt.Printf("Adding project: %v (%v)", project, abbreviation)

	if description != "" {
		fmt.Printf(" - %v\n", description)
	}

	fmt.Printf("\n")

	if err := getDatabase().AddProject(project, abbreviation, description); err != nil {
		fmt.Fprintf(os.Stderr, "Error occured adding project:\n%v\n", err)
	}
}

func taskListCommand(arguments ...string) {
	var project *progress.Project

	if len(arguments) > 0 {
		fmt.Fprintf(os.Stderr, "usage: @task list\n")
		os.Exit(3)
	}

	tasks, err := getDatabase().Tasks(nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not retreive task list: %v\n", err)
		os.Exit(4)
	}

	for _, task := range tasks {
		if project == nil || project.Identifier != task.Project.Identifier {
			project = task.Project
			fmt.Printf("%v:\n", project.Name)
		}

		fmt.Printf("- %v [%v]\n", task.Summary, task.Identifier)
	}
}

func taskAddCommand(arguments ...string) {
	var err error

	var project progress.Project
	var projectRef string

	if len(arguments) == 0 || len(arguments[0]) == 0 {
		fmt.Fprintf(os.Stderr, "usage: @task add [~project] <>")
		os.Exit(3)
	}

	database := getDatabase()

	// If arguments[0] is ~PRJ where PRJ is a project reference, this task goes
	// into that project instead of the default project.
	if arguments[0][0] == '~' {
		projectRef = arguments[0][1:]
		project, err = database.Project(projectRef)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get project (%v): %v\n", projectRef, err)
			os.Exit(4)
		}
	} else {
		project, err = database.DefaultProject()

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", project.Abbreviation)
			os.Exit(5)
		}
	}

	summary := strings.Join(arguments, " ")
	fmt.Printf("Adding task to %v: %v\n", project.Name, summary)

	if _, err = database.AddTask(project, summary); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create task: %v\n", err)
		os.Exit(6)
	}
}
