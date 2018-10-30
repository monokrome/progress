package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/monokrome/progress"
	"gopkg.in/alecthomas/kingpin.v2"

	// SQL dialects
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	kingpin.CommandLine.HelpFlag.Short('h')
}

func main() {
	var (
		err error

		cli = kingpin.New("progress", "a tool for tracking your progress")

		project = cli.Command("project", "command for managing project")

		projectCreate             = project.Command("create", "create a new project")
		projectCreateAbbreviation = projectCreate.Flag("abbreviation", "abbreviation to use for the project").Short('a').String()
		projectCreateName         = projectCreate.Arg("name", "name of the project").Required().String()

		projectRemove             = project.Command("remove", "remove a project from the database")
		projectRemoveAbbreviation = projectRemove.Flag("abbreviation", "abbreviation to use for the project").Short('a').Required().String()

		projectList = project.Command("list", "list projects in the database")

		task = cli.Command("task", "command for managing task")

		taskActive = task.Command("active", "get the currently active task")

		taskCreate             = task.Command("create", "create a new task")
		taskCreateAbbreviation = taskCreate.Flag("abbreviation", "abbreviation for the project to create the task in").Default("").Short('a').String()
		taskCreateTopic        = taskCreate.Arg("topic", "topic of the newly created task").Required().String()

		taskTag       = task.Command("tag", "create a new task on the current task")
		taskTagDetach = taskTag.Flag("detach", "detaches the tag instead of attaching it").Short('d').Bool()
		taskTagName   = taskTag.Arg("tag", "the name of the tag to modify").Required().String()

		taskList             = task.Command("list", "list tasks")
		taskListAbbreviation = taskList.Flag("abbreviation", "abbreviation for the project to list tasks from").Short('a').String()
	)

	cli.UsageTemplate(kingpin.CompactUsageTemplate)

	options, _, err := progress.NewOptions("progress")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
	}

	database, err := gorm.Open(options.Storage.Backend, options.Storage.Options)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database connection: %v\n", err)
		os.Exit(1)
	}

	defer database.Close()

	progress.EnsureSchema(database)

	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {
	case projectCreate.FullCommand():
		err = CreateProject(database, *projectCreateName, *projectCreateAbbreviation)

	case projectRemove.FullCommand():
		err = RemoveProject(database, *projectRemoveAbbreviation)

	case projectList.FullCommand():
		err = ListProjects(database)

	case taskActive.FullCommand():
		err = TaskActive(database)

	case taskCreate.FullCommand():
		err = CreateTask(database, *taskCreateTopic, *taskCreateAbbreviation)

	case taskList.FullCommand():
		err = ListTasks(database, *taskListAbbreviation)

	case taskTag.FullCommand():
		err = TaskTag(database, *taskTagDetach, *taskTagName)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(10)
	}
}
